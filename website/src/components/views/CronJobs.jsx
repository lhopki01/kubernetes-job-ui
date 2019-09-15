import React from 'react';
import { JobStatusIconLink } from '../JobStatusIcon';
import { NavBar } from '../NavBar';

class CronJobs extends React.Component {
  constructor() {
      super();

      this.state = {
          expandedRows : []
      };
  }

  handleRowClick(rowId) {
    const currentExpandedRows = this.state.expandedRows;
    const isRowCurrentlyExpanded = currentExpandedRows.includes(rowId);

    const newExpandedRows = isRowCurrentlyExpanded ?
                            currentExpandedRows.filter(id => id !== rowId) :
                            currentExpandedRows.concat(rowId);

    this.setState({expandedRows : newExpandedRows});
  }

  render() {
    return (
      <React.Fragment>
        <NavBar {...this.props} />
        <CronJobsTable cronJobs={ this.props.cronJobs } onClick={ (i) => this.handleRowClick(i) } expandedRows={ this.state.expandedRows }></CronJobsTable>
      </React.Fragment>
    )
  }
}

class CronJobsTable extends React.Component {
  render() {
    const rows = this.props.cronJobs.map((item, index) => {
      const id = item.name+item.namespace
      const clickCallback = () => this.props.onClick(id);
      var cronJob = this.props.expandedRows.includes(id) ?
                    <tr>
                      <td colSpan="6">
                        <CronJobInformationPanel { ...item }/>
                        <JobTable { ...item } />
                      </td>
                    </tr> :
                    null;
      return (
        <React.Fragment key={item.name+item.namespace}>
        <tr onClick={ clickCallback } key={item.name+item.namespace}>
          <td><a href={"/namespaces/"+item.namespace+"/cronjobs/"+item.name}>{item.name}</a></td>
          <td>{ item.config.description }</td>
          <td><ReturnFirstJob cronJob={item}></ReturnFirstJob></td>
          <td><ReturnPreviousJobs cronJob={item}></ReturnPreviousJobs></td>
          <td><RunButton cronJob={item}></RunButton></td>
        </tr>
        { cronJob }
        </React.Fragment>
      )
    })
    return (
      <div className="container-fluid">
        <table className="table table-hover">
          <tbody>
            <tr>
              <th>Name</th>
              <th>Description</th>
              <th>Last Run</th>
              <th>Previous Runs</th>
              <th>Run</th>
            </tr>
            {rows}
          </tbody>
        </table>
      </div>
    )
  }
}

function CronJobInformationPanel(props) {
  return (
        <div className="container-fluid">
          <div className="alert alert-secondary">
            <div className="row">
              <div className="col-11">
                <h6>Schedule: {props.schedule}</h6>
                <h6>Namespace: {props.namespace}</h6>
              </div>
            </div>
          </div>
        </div>
  )
}

function JobTable(props) {
  const jobs = (props.jobs != null) ?
    props.jobs.map((job, index) => {
      return (
        <tr key={job.name}>
          <td><a href={props.name+"/jobs/"+job.name}>{job.name}</a></td>
          <td>{job.creationTime}</td>
          <td><JobStatusIconLink cronJob={props} job={job} /></td>
        </tr>
      )
    }) :
    null;
  return (
    <div className="container-fluid">
      <table className="table table-condensed">
        <thead>
          <tr>
            <th>Job</th>
            <th>Creation Time</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {jobs}
        </tbody>
      </table>
    </div>
  )
}

function RunButton(props) {
  return(
    <a href={"/createjob?namespace="+props.cronJob.namespace+"&cronjob="+props.cronJob.name}>
      <svg id="i-play" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="25" height="25" fill="none" stroke="currentcolor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" className="text-centre">
        <path d="M10 2 L10 30 24 16 Z" />
      </svg>
    </a>
  )
}

function ReturnFirstJob(props) {
  if (props.cronJob.jobs != null && props.cronJob.jobs.length > 0) {
    return (
      <JobStatusIconLink cronJob={props.cronJob} job={props.cronJob.jobs[0]} />
    )
  }
  return null
}

function ReturnPreviousJobs(props) {
  if (props.cronJob.jobs != null && props.cronJob.jobs.length > 1) {
    return props.cronJob.jobs.slice(1).map(job => {
      return (
        <JobStatusIconLink key={job.name} job={job} cronJob={props.cronJob} />
      )
    })
  }
  return null
}



export { CronJobs }
