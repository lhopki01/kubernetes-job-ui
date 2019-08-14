import React from 'react';
import JobStatusIcon from '../JobStatusIcon.jsx';

class CronJob extends React.Component {
  render() {
    console.log(this.props)
    console.log("----")
    return (
        <CronJobTable {...this.props}></CronJobTable>
    );
  }
}

export default CronJob

function CronJobTable(props) {
    const jobs = props.cronJobs.map((item, index) => {
        console.log(item.Name+"="+props.match.params.cronJobName)
        if (item.Name === props.match.params.cronJobName && item.Namespace === props.match.params.namespace) {
            console.log("in logic")
            return (item.Jobs.map((job, index) => {
                return (
                    <tr key={job.Name}>
                        <td><a href={item.Name+"/jobs/"+job.Name}>{item.Name}</a>{job.Name}</td>
                        <td>{job.Namespace}</td>
                        <td>{job.CreationTime}</td>
                        <td><JobStatusIcon CronJob={item} Job={job} /></td>
                    </tr>
                )
            }))
        }
        return null
    })
    return (
        <table>
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
    )
}
