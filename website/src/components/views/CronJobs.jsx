import React from 'react';
import { JobStatusIconLink } from '../JobStatusIcon';
import { NavBar } from '../NavBar';

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
                <tr key={item.name+item.namespace}>
                    <td><a href={"namespaces/"+item.namespace+"/cronjobs/"+item.name}>{item.name}</a></td>
                    <td>{item.namespace}</td>
                    <td>{item.schedule}</td>
                    <td><ReturnFirstJob cronJob={item}></ReturnFirstJob></td>
                    <td><ReturnPreviousJobs cronJob={item}></ReturnPreviousJobs></td>
                    <td><RunButton cronJob={item}></RunButton></td>
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
