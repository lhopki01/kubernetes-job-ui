import React from 'react';

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
            console.log(url)
            fetch(url)
            .then((reponse) =>
                reponse.json()
            )
            .then((jsonData) => {
                console.log(jsonData[0].Phase)
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
            console.log("Clearing interval")
            clearInterval(this.interval)
        }
    } 

    componentDidMount() {
        this.getLogs()
        this.interval = setInterval(() => this.getLogs(), 1000);
    }

    render() {
        console.log(this.props)
        if (this.state.isLoading) {
            return (
                <div id="job">
                  This page will have a Job on it.
                </div>
            )
        } else {
            return (
                <LogsTable job={this.state.job} />
            )
        }
    }
}


function LogsTable(props) {
    const logs = props.job.map((item, index) => {
        return (item.Containers.map((c, index) => {
            return (
                <>
                    <tr>
                        <td>{item.Name}</td>
                        <td>{item.CreationTime}</td>
                        <td>{item.Phase}</td>
                    </tr>
                    <tr><th>Container: {c.Name}</th></tr>
                    <tr><td><pre>{c.Logs}</pre></td></tr>
                </>
            )
        }))
    })
    return (
        <table>
            <thead>
                <tr>
                    <th>Pod</th>
                    <th>Creation Time</th>
                    <th>Phase</th>
                </tr>
            </thead>
            <tbody>
                {logs}
            </tbody>
        </table>
    )

}

export default Job

