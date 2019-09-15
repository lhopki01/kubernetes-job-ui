import React from 'react';
import { JobStatusIcon } from '../JobStatusIcon';
import { NavBar } from '../NavBar';

class Job extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      isLoading: true,
      shouldPoll: true,
    }
  }

  shouldPoll(status) {
    if (status === "active") {
      return true
    }
    return false
  }
  isLoading(jsonData) {
    if (jsonData.pods !== null) {
      return false
    }
    return true
  }

  async getLogs() {
    if (this.state.shouldPoll) {
      try {
        const {namespace, cronJobName, jobName} = this.props.match.params
        const url = `/api/v1/namespaces/${namespace}/cronjobs/${cronJobName}/jobs/${jobName}`
        const response = await fetch(url)
        const jsonData = await response.json()
        this.setState({
          job: jsonData,
          shouldPoll: this.shouldPoll(jsonData.status),
          isLoading: this.isLoading(jsonData)
        })
        return

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
          <div id="job" className="container-fluid">
            <pre>Loading...</pre>
          </div>
        </React.Fragment>
      )
    }
    return (
      <React.Fragment>
        <NavBar {...this.props} />
        <div className="container-fluid">
          <JobInformationPanel job={this.state.job} />
          <PodTabs job={this.state.job} />
        </div>
      </React.Fragment>
    )
  }
}

function JobInformationPanel(props) {
  return (
    <div className="alert alert-secondary">
      <h4>{props.job.name} <JobStatusIcon status={props.job.status} /></h4>
      <h6>Namespace: {props.job.namespace}</h6>
      <h6>Creation Time: {props.job.creationTime}</h6>
    </div>
  )
}

function PodTabs(props) {
  const tabs = (props.job.pods.map((p, index) => {
    let active = ""
    if (index === 0) {
      active="active"
    }
    return (
      <a key={p.name} className={"nav-item nav-link "+active} data-toggle="tab" href={"#"+p.name} role="tab"><h5>{p.name} <JobStatusIcon status={p.status}/></h5></a>
    )
  }))
  const tabContent = (props.job.pods.map((p, index) => {
    let active = ""
    if (index === 0) {
      active="active show"
    }
    return (
      <div key={p.name} className={"tab-pane fade "+active} id={p.name} role="tabpanel">
        <ContainerTabs containers={p.containers} jobName={p.name}/>
      </div>
    )
  }))
  return (
    <div className="container-fluid">
    <React.Fragment>
    <div className="nav nav-tabs" id="nav-tab" role="tablist">
      <div className="navbar-brand">Pods:</div>
      {tabs}
    </div>
    <div className="tab-content" id="nav-tabContent" role="tabpanel">
      {tabContent}
    </div>
    </React.Fragment>
    </div>
  )
}

function ContainerTabs(props) {
  const tabs = (props.containers.map((c, index) => {
    let active = ""
    if (index === 0) {
      active="active"
    }
    return (
      <a key={c.name} className={"nav-item nav-link "+active} data-toggle="tab" href={"#"+props.jobName+c.name} role="tab">{c.name}</a>
    )
  }))
  const tabContent = (props.containers.map((c, index) => {
    let active = ""
    if (index === 0) {
      active="active show"
    }
    return (
      <div key={c.name} className={"tab-pane fade "+active} id={props.jobName+c.name} role="tabpanel">
        <pre className="logs">{c.logs}</pre>
      </div>
    )
  }))
  return (
    <React.Fragment>
    <div className="nav nav-tabs" id="nav-tab" role="tablist">
      <div className="navbar-brand">Containers:</div>
      {tabs}
    </div>
    <div className="tab-content" id="nav-tabContent">
      {tabContent}
    </div>
    </React.Fragment>
  )

}

export { Job }

