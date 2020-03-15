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
      {/* <CronJobsTable cronJobs={ this.props.cronJobs } onClick={ (i) => this.handleRowClick(i) } expandedRows={ this.state.expandedRows }></CronJobsTable> */}
      </React.Fragment>
    )
  }
}

function Folders(props) {
  console.log(props.expandedFolders)
  console.log(props.expandedRows)
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
    const clickCallback = () => props.folderClick(folder);
    var cronJobs = props.expandedFolders.includes(folder) ?
                   <FolderCronJobsTable folder={folder} cronJobs={props.cronJobs} rowClick={props.rowClick} expandedRows={props.expandedRows}/> :
                   null;
    var symbol = "  +"
    if (props.expandedFolders.includes(folder)) {
      var symbol = "  -"
    }
    return (
      <React.Fragment key={folder}>
      <tr onClick={ clickCallback } key={folder}>
        <td><h4>{folder}{symbol}</h4></td>
      </tr>
      { cronJobs }
      </React.Fragment>
    )
  })

  return(
    <div className="container-fluid">
      <table className="table table-hover">
        <tbody>
          {folderRows}
          <FolderCronJobsTable folder="" cronJobs={props.cronJobs} rowClick={props.rowClick} expandedRows={props.expandedRows}/>
        </tbody>
      </table>
    </div>
  )
}

function FolderCronJobsTable(props) {
  //console.log(props.folder)
  //console.log(props.cronJobs)
  console.log(props.expandedRows)
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
          <td>{item.name}</td>
          <td>{ item.config.shortDescription }</td>
          <td className="text-center"><ReturnFirstJob cronJob={item}></ReturnFirstJob></td>
          <td><ReturnPreviousJobs cronJob={item}></ReturnPreviousJobs></td>
          <td className="text-center"><RunButton cronJob={item}></RunButton></td>
        </tr>
        {cronJob}
        </React.Fragment>
      )
    }
    return null
  console.log(cronJobs)
  })
  return (
    <div className="container-fluid">
          <div className="alert alert-light">
      <table className="table">
        <tbody>
          <tr>
            <th>Name</th>
            <th>Description</th>
            <th className="text-center">Last Run</th>
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

class CronJobsTable extends React.Component {
  render() {
    const rows = this.props.cronJobs.map((item, index) => {
      console.log(item.config.folder)
      const id = item.name+item.namespace
      const clickCallback = () => this.props.onClick(id);
      var cronJob = this.props.expandedRows.includes(id) ?
                    <tr>
                      <td colSpan="6">
                        <CronJobInformationPanel { ...item }/>
                      </td>
                    </tr> :
                    null;
      return (
        <React.Fragment key={item.name+item.namespace}>
        <tr onClick={ clickCallback } key={item.name+item.namespace}>
          <td><a href={"/namespaces/"+item.namespace+"/cronjobs/"+item.name}>{item.name}</a></td>
          <td>{ item.config.shortDescription }</td>
          <td className="text-center"><ReturnFirstJob cronJob={item}></ReturnFirstJob></td>
          <td><ReturnPreviousJobs cronJob={item}></ReturnPreviousJobs></td>
          <td className="text-center"><RunButton cronJob={item}></RunButton></td>
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
              <th className="text-center">Last Run</th>
              <th>Previous Runs</th>
              <th className="text-center">Links</th>
            </tr>
            {rows}
          </tbody>
        </table>
      </div>
    )
  }
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
        <tr key={job.name} >
          <td>{job.name}</td>
          <td>{job.creationTime}</td>
          <td className="text-center"><JobStatusIconLink cronJob={props} job={job} /></td>
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



export { CronJobs }
