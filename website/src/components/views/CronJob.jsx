import React from 'react';
import { JobStatusIconLink } from '../JobStatusIcon';
import { NavBar } from '../NavBar';

class CronJob extends React.Component {
  render() {
    return (
      <React.Fragment>
        <NavBar {...this.props} />
        <CronJobTable {...this.props}></CronJobTable>
      </React.Fragment>
    );
  }
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
