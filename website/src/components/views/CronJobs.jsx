import React from 'react';
import JobStatusIcon from '../JobStatusIcon.jsx';
import NavBar from '../NavBar.jsx';

class CronJobs extends React.Component {
    render() {
        return (
            <React.Fragment>
                <NavBar {...this.props} />
                <CronJobsTable {...this.props}></CronJobsTable>
            </React.Fragment>
        )
    }
}

class CronJobsTable extends React.Component {

    render() {
        const rows = this.props.cronJobs.map((item, index) => {
            return (
                <tr key={item.Name+item.Namespace}>
                    <td><a href={"namespaces/"+item.Namespace+"/cronjobs/"+item.Name}>{item.Name}</a></td>
                    <td>{item.Namespace}</td>
                    <td>{item.Schedule}</td>
                    <td><ReturnFirstJob CronJob={item}></ReturnFirstJob></td>
                    <td><ReturnPreviousJobs CronJob={item}></ReturnPreviousJobs></td>
                    <td><RunButton CronJob={item}></RunButton></td>
                </tr>
            )
        })
        return (
            <div className="container-fluid">
                <table className="table table-condensed table-bordered table-striped">
                    <tbody>
                        <tr>
                            <th>Name</th>
                            <th>Namespace</th>
                            <th>Schedule</th>
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

function RunButton(props) {
    return(
        <a href={"/createjob?namespace="+props.CronJob.Namespace+"&cronjob="+props.CronJob.Name}>
            <svg id="i-play" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="25" height="25" fill="none" stroke="currentcolor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" className="text-centre">
                <path d="M10 2 L10 30 24 16 Z" />
            </svg>
        </a>
    )
}

function ReturnFirstJob(props) {
    if (props.CronJob.Jobs != null && props.CronJob.Jobs.length > 0) {
        return (
            <JobStatusIcon CronJob={props.CronJob} Job={props.CronJob.Jobs[0]} />
        )
    }
    return (
        null
    )
}

function ReturnPreviousJobs(props) {
    if (props.CronJob.Jobs != null && props.CronJob.Jobs.length > 1) {
        return props.CronJob.Jobs.slice(1).map(job => {
            return (
                <JobStatusIcon key={job.Name} Job={job} CronJob={props.CronJob} />
            )
        })
    }
    return (
        null
    )
}



export default CronJobs;
