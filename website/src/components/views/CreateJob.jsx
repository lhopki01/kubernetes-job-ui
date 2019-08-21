import React from 'react';
import queryString from 'query-string'
import { NavBar } from '../NavBar';

class CreateJob extends React.Component {
    constructor(props) {
        super(props);
        this.onSubmitWithProps = this.onSubmitWithProps.bind(this)

        const queryValues = queryString.parse(props.location.search)
        let cronJob = {}
        props.cronJobs.map(item => {
            if (queryValues.namespace === item.namespace && queryValues.cronjob === item.name) {
                cronJob = item
            }
            return null
        })
        let formValues = {}
        cronJob.config.options.map((option, index) => {
            formValues[index] = option.default
            return null
        })

        this.state = {
            errors: {},
            queryValues: queryValues,
            cronJob: cronJob,
            formValues,
        }
    }

    async onSubmitWithProps(event, props) {
        event.preventDefault();
        let jobRequest = new Array(props.cronJob.config.options.length)
        let formValues = {}
        props.cronJob.config.options.map((option, index) => {
            jobRequest[index] = {
                "envVar": option.envVar,
                "container": option.container,
                "value": event.target[index].value,
            }
            formValues[index] = event.target[index].value
            return null
        })
        this.setState({
            formValues: formValues
        })
        const url = "api/v1/namespaces/"+props.cronJob.namespace+"/cronjobs/"+props.cronJob.name
        try {
            const response = await fetch(url, {
                method: 'post',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(jobRequest),
            })
            const jsonData = await response.clone().json()
            if (response.status === 200) {
                console.log(jsonData)
                this.props.history.push("/namespaces/"+props.cronJob.namespace+"/cronjobs/"+props.cronJob.name+"/jobs/"+jsonData.job)
                return
            }
            let errors = {}
            jsonData.map(error => {
                errors[error.optionIndex] = error.error
                return null
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
    const option = (
        <div className="form-group row">
            <label className="col text-right col-form-label">
                {props.option.envVar}
            </label>
            <div className="col">
                <input
                    className="form-control"
                    placeholder={props.option.default}
                    defaultValue={props.formValues[props.index]}
                    name={props.index}
                />
            </div>
            <div className="col col-form-label">
                <span>{props.option.description}</span>
            </div>
        </div>
    )
    if (props.index in props.errors) {
        return (
            <div>
            {option}
            <div className="form-group row">
                <label className="col text-right col-form-label">
                </label>
                <div className="col">
                    <span style={{color: "red"}}>{props.errors[props.index]}</span>
                </div>
                <div className="col col-form-label">
                </div>
            </div>
            </div>
        )
    }
    return (
        option
    )
}

export { CreateJob };
