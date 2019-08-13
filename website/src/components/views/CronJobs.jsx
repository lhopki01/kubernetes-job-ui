import React from 'react';

class CronJobs extends React.Component {
    render() {
        console.log(this.props.str)
        console.log(this.props.cronJobs)
        //return (<span>Cronjobs</span>)
        return (
            //<CronJobTable cronJobs={this.props.cronJobs}></CronJobTable>
            <CronJobsTable {...this.props}></CronJobsTable>
        )
    }
}

class CronJobsTable extends React.Component {

    render() {
        const rows = this.props.cronJobs.map((item, index) => {
            return (
                <tr key={item.Name+item.Namespace}>
                    <td><a href={"/cronjob?cronjob="+item.Name}>{item.Name}</a></td>
                    <td>{item.Namespace}</td>
                    <td>{item.Schedule}</td>
                    <td><ReturnFirstJob cronJob={item}></ReturnFirstJob></td>
                    <td><ReturnPreviousJobs cronJob={item}></ReturnPreviousJobs></td>
                    <td><RunButton cronJob={item}></RunButton></td>
                </tr>
            )
        })
        return (
            <table>
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
        )
    }
}

function RunButton(props) {
    return(
        <a href={"/createjob?cronjob="+props.cronJob.Name}>
            <svg id="i-play" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="25" height="25" fill="none" stroke="currentcolor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" className="text-centre">
                <path d="M10 2 L10 30 24 16 Z" />
            </svg>
        </a>
    )
}

function ReturnFirstJob(props) {
    if (props.cronJob.Jobs != null && props.cronJob.Jobs.length > 0) {
        return (
            <a href={"/job?cronJob="+props.cronJob.Name+"&job="+props.cronJob.Jobs[0].Name} data-toggle="tooltip" data-placement="bottom" data-original-title={ props.cronJob.Jobs[0].CreationTime }>
                <JobStatusIcon Status={props.cronJob.Jobs[0].Status} Manual={props.cronJob.Jobs[0].Manual} />
            </a>
        )
    }
    return (
        null
    )
}

function ReturnPreviousJobs(props) {
    if (props.cronJob.Jobs != null && props.cronJob.Jobs.length > 1) {
        return props.cronJob.Jobs.slice(1).map(job => {
            return (
                <a key={props.cronJob.Namespace+"-"+job.Name} href={"/job?cronJob="+props.cronJob.Name+"&job="+job.Name} data-toggle="tooltip" data-placement="bottom" data-original-title={ job.CreationTime }>
                    <JobStatusIcon Status={job.Status} Manual={job.Manual} />
                </a>
            )
        })
    }
    return (
        null
    )
}

function JobStatusIcon(props) {
    return (
        <svg viewBox="0 0 32 32" width="20" height="20" fill="none" strokeLinecap="round" strokeLinejoin="round" strokeWidth="3">
        {(() => {
            switch(props.Status) {
                case "succeeded": return <path d="M2 20 L12 28 30 4" stroke="green"/>
                case "failed": return <path d="M2 30 L30 2 M30 30 L2 2" stroke="red"/>
                case "active": return <circle cx="16" cy="16" r="12" stroke="orange"/>
                default: return <path d="M2 30 L30 2 M30 30 L2 2" stroke="black"/>
            }
        })()}
        { props.Manual ? <path d="M19,31 l0,-10 m0,2 q3,-3 6,0 l 0,8 m0,-8 q3,-3 6,0 l 0,8" stroke="black" strokeWidth="2"/> : null}
        </svg>
    )
}

export default CronJobs;
