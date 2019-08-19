import React from 'react';
import JobStatusIcon from '../JobStatusIcon.jsx';
import NavBar from '../NavBar.jsx';

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

export default CronJob

function CronJobTable(props) {
    const jobs = props.cronJobs.map((item, index) => {
        if (item.Name === props.match.params.cronJobName && item.Namespace === props.match.params.namespace) {
            return (item.Jobs.map((job, index) => {
                return (
                    <tr key={job.Name}>
                        <td><a href={item.Name+"/jobs/"+job.Name}>{job.Name}</a></td>
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
