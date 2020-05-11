import React from 'react';
import queryString from 'query-string'
import { NavBar } from '../NavBar';

class CreateJob extends React.Component {
  constructor(props) {
    super(props);

    const queryValues = queryString.parse(props.location.search)
    const cronJob = props.cronJobs.reduce((matchedCronJob, item) => {
      if (queryValues.namespace === item.namespace && queryValues.cronjob === item.name) {
        matchedCronJob = item
      }
      return matchedCronJob
    }, {})
    let formValues = []
    if (cronJob.config.options !== null) {
      formValues = cronJob.config.options.map((option, index) => {
        return option.default
      })
    }

    this.state = {
      errors: {},
      queryValues: queryValues,
      cronJob: cronJob,
      formValues,
    }
  }

  onSubmitWithProps = async (event, props) => {
    event.preventDefault();
    const options = props.cronJob.config.options
    const jobRequest = options.map((option, index) => {
      return {
        "envVar": option.envVar,
        "container": option.container,
        "value": event.target[index].value,
      }
    })
    let formValues = options.map((option, index) => {
      return event.target[index].value
    })
    this.setState({
      formValues: formValues
    })
    const { namespace, name } = props.cronJob
    const url = `/api/v1/namespaces/${namespace}/cronjobs/${name}`
    try {
      const response = await fetch(url, {
        method: 'post',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(jobRequest),
      })
      const jsonData = await response.json()
      if (response.status === 200) {
        const { namespace, name } = props.cronJob
        this.props.history.push(`/namespaces/${namespace}/cronjobs/${name}/jobs/${jsonData.job}`)
        return
      }
      // using an array because want to to lookup by number
      let errors = {}
      jsonData.forEach(error => {
        errors[error.optionIndex] = error.error
      })
      this.setState({
        errors: errors
      })
    }
    catch(error) {
      console.error(error)
    }
  }

  render() {
    return(
      <React.Fragment>
        <NavBar {...this.props} />
        <JobForm onSubmitWithProps={this.onSubmitWithProps} {...this.state}/>
      </React.Fragment>
    )
  }
}


function JobForm(props) {
  const onSubmit = (event) => {
      props.onSubmitWithProps(event, props)
  }
  if (props.cronJob.config.errors !== null) {
    const errors = props.cronJob.config.errors.map((error, index) => {
      return (
        <pre key={index} className="linewrap">{error}</pre>
      )
    })
    return(
      <React.Fragment>
      <div className="container">
        <div className="alert alert-secondary">
          {errors}
        </div>
        <pre>{props.cronJob.config.raw}</pre>
      </div>
      </React.Fragment>
    )
  }
  return(
    <React.Fragment key={props.cronJob.name}>
      <div className="container">
        <form onSubmit={onSubmit}>
          <Options options={props.cronJob.config.options} errors={props.errors} {...props}/>
          <div className="form-group row">
            <div className="col"></div>
            <div className="col">
              <button type="submit" className="btn btn-primary float-right">Run</button>
            </div>
            <div className="col"></div>
          </div>
        </form>
      </div>
    </React.Fragment>
  )
}

function Options(props) {
  let containerName = ""
  return ( props.cronJob.config.options.map((option, index) => {
    if (containerName !== option.container) {
      containerName = option.container
      return (
        <React.Fragment key={index}>
          <div className="row">
            <div className="col" />
            <div className="col">
              <h4>Container: {option.container}</h4>
            </div>
            <div className="col" />
          </div>
          <Option option={option} index={index} {...props}/>
        </React.Fragment>
      )
    } else {
      return (
        <Option key={index} option={option} index={index} {...props}/>
      )
    }
  }))
}

function Option(props) {
  return (
    <div className="form-group row">
      <label className="col text-right col-form-label">
        {props.option.envVar}
      </label>
      <div className="col">
        <FormInput option={props.option} defaultValue={props.formValues[props.index]} name={props.index} />
        <span style={{color: "red"}}>{props.errors[props.index]}</span>
      </div>
      <div className="col col-form-label">
        <span>{props.option.description}</span>
      </div>
    </div>
  )
}

function FormInput(props) {
  if (props.option.type === "list" || props.option.type === "bool") {
    const selectOptions = props.option.values.map(selectOption => {
      return (
        <option key={selectOption}>{selectOption}</option>
      )
    })
    return(
      <select className="form-control" name={props.name}>
        {selectOptions}
      </select>

    )
  }
  if (props.option.type === "textarea") {
  return(
    <textarea
      className="form-control"
      placeholder={props.option.default}
      defaultValue={props.defaultValue}
      name={props.name}
    />
  )

  }
  return(
    <input
      className="form-control"
      placeholder={props.option.default}
      defaultValue={props.defaultValue}
      name={props.name}
    />
  )

}

export { CreateJob };
