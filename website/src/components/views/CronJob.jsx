import React from 'react';
import { JobStatusIconLink } from '../JobStatusIcon';
import { NavBar } from '../NavBar';

class CronJob extends React.Component {
  constructor(props) {
    super(props);
    props.cronJobs.forEach(item => {
      if (item.name === props.match.params.cronJobName && item.namespace === props.match.params.namespace) {
        this.state = {
          cronJob: item
       }
      }
    })
  }
  render() {
    return (
      <React.Fragment>
        <NavBar {...this.props} />
        <CronJobInformationPanel { ...this.state.cronJob } />
        <CronJobTable {...this.props}></CronJobTable>
      </React.Fragment>
    );
  }
}

function CronJobInformationPanel(props) {
  console.log(props)
  return (
    <div className="container-fluid">
      <div className="alert alert-secondary">
        <div className="row">
          <div className="col-11">
            <h4>{props.name}</h4>
            <h6>{props.config.description}</h6>
            <h6>Namespace: {props.namespace}</h6>
            <h6>Schedule: {props.schedule}</h6>
          </div>
          <div className="col align-middle">
            <RunButton cronJob={props} />
          </div>
        </div>
      </div>
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

function CronJobTable(props) {
  const jobs = props.cronJobs.map((item, index) => {
    if (item.name === props.match.params.cronJobName && item.namespace === props.match.params.namespace) {
      return (item.jobs.map((job, index) => {
        return (
          <tr key={job.name}>
            <td><a href={item.name+"/jobs/"+job.name}>{job.name}</a></td>
            <td>{job.namespace}</td>
            <td>{job.creationTime}</td>
            <td><JobStatusIconLink cronJob={item} job={job} /></td>
          </tr>
        )
      }))
    }
    return null
  })
  return (
    <div className="container-fluid">
      <table className="table table-condensed table-bordered table-striped">
        <thead>
          <tr>
            <th>Name</th>
            <th>Namespace</th>
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

export { CronJob }
