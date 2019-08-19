import React from 'react';
import NavBar from '../NavBar.jsx';

class Job extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            isLoading: true,
            poll: true,
        }
    }

    getLogs() {

        if (this.state.poll) {
            const url = "http://localhost:8080/api/v1/namespaces/"+this.props.match.params.namespace+"/jobs/"+this.props.match.params.jobName;
            fetch(url)
            .then((response) =>
                response.json()
            )
            .then((jsonData) => {
                if (jsonData[0].Phase !== "Running") {
                    this.setState({
                        poll: false,
                    })
                }
                this.setState({
                    job: jsonData,
                    isLoading: false,
                })
                return jsonData
            })
            .catch((error) => {
                console.error(error)
            })
        } else {
            clearInterval(this.interval)
        }
    } 

    componentDidMount() {
        this.getLogs()
        this.interval = setInterval(() => this.getLogs(), 1000);
    }

    render() {
        if (this.state.isLoading) {
            return (
                <React.Fragment>
                    <NavBar {...this.props} />
                    <div id="job">
                      Loading...
                    </div>
                </React.Fragment>
            )
        } else {
            return (
                <React.Fragment>
                    <NavBar {...this.props} />
                    <LogsTable job={this.state.job} />
                </React.Fragment>
            )
        }
    }
}


function ContainerLogs(props) {
    return (props.Containers.map((c, index) => {
        return (
            <React.Fragment key={c.Name}>
                <tr><th colSpan="4">Container: {c.Name}</th></tr>
                <tr><td colSpan="4"><pre>{c.Logs}</pre></td></tr>
            </React.Fragment>
        )
    }))
}


function LogsTable(props) {
    const logs = props.job.map((item, index) => {
        return (
            <React.Fragment key={index}>
                <tr>
                    <th>Pod</th>
                    <th>Creation Time</th>
                    <th>Phase</th>
                </tr>
                <tr>
                    <td>{item.Name}</td>
                    <td>{item.CreationTime}</td>
                    <td>{item.Phase}</td>
                </tr>
                <ContainerLogs Containers={item.Containers} PodName={item.Name} />
            </React.Fragment>
        )
    })
    return (
        <div className="container-fluid">
            <table className="table table-condensed table-bordered table-striped">
                <tbody>
                    {logs}
                </tbody>
            </table>
        </div>
    )

}

export default Job

