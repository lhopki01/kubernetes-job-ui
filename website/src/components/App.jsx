import React from 'react';
import { BrowserRouter as Router, Route, Redirect } from "react-router-dom";
import '../stylesheets/App.css';
import { CronJobs } from './views/CronJobs';
import { CronJob } from './views/CronJob';
import { Job } from './views/Job';
import { CreateJob } from './views/CreateJob';

class App extends React.Component {
    constructor(props) {
      super(props);
      this.state = {
        isLoading: true,
      };
    }

    async componentDidMount() {
        try {
            const response = await fetch('/api/v1/cronjobs')
            const jsonData = await response.json()
            console.log(jsonData)
            this.setState({
               cronJobs: jsonData,
               isLoading: false,
            })
        } catch(error) {
            console.error(error)
        }
    }

    render() {
        const { isLoading, cronJobs } = this.state
        if (isLoading){
            return (
                <pre>Loading</pre>
            )
        }
        return (
            <Router>
                <Route exact path="/" render={() => <Redirect to='/cronjobs' />}></Route>
                <Route path="/cronjobs" render={() => <CronJobs cronJobs={cronJobs}/>}></Route>
                <Route
                    exact path="/namespaces/:namespace/cronjobs/:cronJobName"
                    render={
                        props => <CronJob
                            {...props}
                            cronJobs={cronJobs}
                        />}>
                </Route>
                <Route path="/namespaces/:namespace/cronjobs/:cronJobName/jobs/:jobName" component={Job}></Route>
                <Route path="/createjob" render={props => <CreateJob {...props} cronJobs={cronJobs}/>}></Route>
            </Router>
        );
    }
}

export { App }
