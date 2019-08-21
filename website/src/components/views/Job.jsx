import React from 'react';
import { NavBar } from '../NavBar';

class Job extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            isLoading: true,
            poll: true,
        }
    }

    async getLogs() {
        if (this.state.poll) {
            try {
                const url = "/api/v1/namespaces/"+this.props.match.params.namespace+"/jobs/"+this.props.match.params.jobName;
                const response = await fetch(url)
                const jsonData = await response.json()
                if (jsonData[0].phase !== "Running") {
                    this.setState({
                        poll: false,
                    })
                }
                this.setState({
                    job: jsonData,
                    isLoading: false,
                })
                return jsonData

            } catch(error) {
                console.error(error)
            }
        } else {
            clearInterval(this.interval)
        }
    }

    componentDidMount() {
        this.getLogs()
        this.interval = setInterval(() => this.getLogs(), 2000);
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
    return (props.containers.map((c, index) => {
        return (
            <React.Fragment key={c.name}>
                <tr><th colSpan="4">Container: {c.Name}</th></tr>
                <tr><td colSpan="4"><pre>{c.logs}</pre></td></tr>
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
                    <td>{item.name}</td>
                    <td>{item.creationTime}</td>
                    <td>{item.phase}</td>
                </tr>
                <ContainerLogs containers={item.containers} />
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

export { Job }

