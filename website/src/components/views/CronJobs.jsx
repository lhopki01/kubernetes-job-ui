import React from 'react';
import { JobStatusIconLink } from '../JobStatusIcon';
import { NavBar } from '../NavBar';

class CronJobs extends React.Component {
  constructor() {
    super();
    this.state = {
      expandedRows : [],
      expandedFolders : []
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
  handleFolderClick(rowId) {
    const currentExpandedFolders = this.state.expandedFolders;
    const isFolderCurrentlyExpanded = currentExpandedFolders.includes(rowId);
    const newExpandedFolders = isFolderCurrentlyExpanded ?
                            currentExpandedFolders.filter(id => id !== rowId) :
                            currentExpandedFolders.concat(rowId);
    this.setState({expandedFolders : newExpandedFolders});
  }

  render() {
    return (
      <React.Fragment>
        <NavBar {...this.props} />
        <Folders
          cronJobs={ this.props.cronJobs }
          folderClick={ (i) => this.handleFolderClick(i) }
          rowClick={ (i) => this.handleRowClick(i) }
          expandedFolders={ this.state.expandedFolders }
          expandedRows={ this.state.expandedRows }>
        </Folders>
      </React.Fragment>
    )
  }
}

function Folders(props) {
  if (props.cronJobs === null) {
    return null
  }
  let folders = props.cronJobs.map(item => {
    return item.config.folder
  })
  folders = [...new Set(folders)].sort()
  const folderRows = folders.map(folder => {
    if (folder === "") {
      return null
    }
    return (
      <FolderRow folder={folder} {...props}/>
    )
  })

  return(
    <div className="container-fluid">
      <table className="table table-hover">
        <tbody>
          {folderRows}
          <FolderRow folder="" {...props}/>
        </tbody>
      </table>
    </div>
  )
}

function FolderRow(props) {
  if (props.cronJobs === null) {
    return null
  }
  const clickCallback = () => props.folderClick(props.folder);
  var cronJobs = props.expandedFolders.includes(props.folder) ? null :
                 <FolderCronJobsTable folder={props.folder} cronJobs={props.cronJobs} rowClick={props.rowClick} expandedRows={props.expandedRows}/>;
  var symbol = DownCheveron()
  if (props.expandedFolders.includes(props.folder)) {
    symbol = RightCheveron()
  }
  return (
    <React.Fragment key={props.folder}>
        <tr onClick={ clickCallback } key={props.folder}>
          <td><h4>{symbol}   {props.folder}</h4></td>
        </tr>
        { cronJobs }
    </React.Fragment>
  )
}

function FolderCronJobsTable(props) {
  const cronJobs = props.cronJobs.map((item) => {
    if (item.config.folder === props.folder) {
      const id = item.name+item.namespace
      const clickCallback = () => props.rowClick(id);
      var cronJob = props.expandedRows.includes(id) ?
                    <tr>
                      <td colSpan="6">
                        <CronJobInformationPanel { ...item }/>
                      </td>
                    </tr> :
                    null;
      return (
        <React.Fragment key={item.name+item.namespace}>
        <tr onClick={clickCallback} key={item.name+item.namespace}>
          <td width="15%" nowrap="nowrap">{item.name}</td>
          <td width="40%">{ item.config.shortDescription }</td>
          <td width="5%" className="text-center"><ReturnFirstJob cronJob={item}></ReturnFirstJob></td>
          <td width="40%"><ReturnPreviousJobs cronJob={item}></ReturnPreviousJobs></td>
          <td width="5%" className="text-center" nowrap="nowrap"><RunButton cronJob={item}></RunButton></td>
        </tr>
        {cronJob}
        </React.Fragment>
      )
    }
    return null
  })
  return (
    <div className="container-fluid">
      <div className="alert alert-light">
        <table className="table">
          <tbody>
            <tr>
              <th>Name</th>
              <th>Description</th>
              <th className="text-center" nowrap="nowrap">Last Run</th>
              <th>Previous Runs</th>
              <th className="text-center">Links</th>
            </tr>
            {cronJobs}
          </tbody>
        </table>
      </div>
    </div>
  )

}

function CronJobInformationPanel(props) {
  const description = props.config.longDescription !== "" ? props.config.longDescription : props.config.shortDescription
  return (
        <div className="container-fluid">
          <div className="alert alert-secondary">
            <div className="row">
              <div className="col-3">
                <h6>Schedule: {props.schedule}</h6>
                <h6>Namespace: {props.namespace}</h6>
              </div>
              <div className="col-8">
                <p>{ description }</p>
              </div>
            </div>
            <div className="row">
              <JobTable {...props} />
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
          <td>{job.name}</td>
          <td>{job.creationTime}</td>
          <td className="text-center"><JobStatusIconLink cronJob={props} job={job} /></td>
          <td><a href={"/namespaces/"+job.namespace+"/cronJobs/"+props.name+"/jobs/"+job.name}>{LogsIcon()}Logs</a></td>
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
            <th className="text-center">Status</th>
            <th>Links</th>
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
      Run
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

function DownCheveron() {
  return (
    <svg viewBox="0 0 32 32" width="18" height="18" fill="none" stroke-width="8">
      <path d="M2 4 L16 20 30 4" stroke="white" />
    </svg>
  )
}

function RightCheveron() {
  return (
    <svg viewBox="0 0 32 32" width="18" height="18" fill="none" stroke-width="8">
      <path d="M4 2 L20 16 4 30" stroke="white" />
    </svg>
  )
}

function LogsIcon() {
  return (
    <svg viewBox="0 0 32 32" width="24" height="24" fill="none" >
      <path d="M6 2 L6 30 26 30 26 10 18 2 Z M18 2 L18 10 26 10" stroke="white" stroke-width="3"/>
    </svg>
  )
}

export { CronJobs }
