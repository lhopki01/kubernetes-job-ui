import React from 'react';
import { BrowserRouter as Router, Route } from "react-router-dom";
import '../stylesheets/App.css';
import CronJobs from './views/CronJobs.jsx';
import CronJob from './views/CronJob.jsx';
import Job from './views/Job.jsx';
import CreateJob from './views/CreateJob.jsx';

class App extends React.Component {
    constructor(props) {
      super(props);
      this.state = {
        isLoading: true,
      };
    }

    componentDidMount() {
       fetch('http://localhost:8080/api/v1/cronjobs')
        .then((response) => 
            response.json()
        )
        .then((jsonData) => {
            this.setState({
                cronJobs: jsonData,
                isLoading: false,
                str: "foobar",
            })
        })
        .catch((error) => {
            console.error(error)
        })
    }
    render() {
        if (this.state.isLoading){
            return (
                <pre>Loading</pre>
            )
        } else {
            return (
                <Router>
                    <Route exact path="/" render={() => <CronJobs cronJobs={this.state.cronJobs}/>}></Route>
                    <Route path="/cronjobs" render={() => <CronJobs cronJobs={this.state.cronJobs}/>}></Route>
                    <Route
                        exact path="/namespaces/:namespace/cronjobs/:cronJobName"
                        render={
                            props => <CronJob
                                {...props}
                                cronJobs={this.state.cronJobs}
                                //cronJobName={props.match.params.cronJobName}
                            />}>
                    </Route>
                    <Route path="/namespaces/:namespace/cronjobs/:cronJobName/jobs/:jobName" component={Job}></Route>
                    <Route path="/createjob" render={props => <CreateJob {...props} cronJobs={this.state.cronJobs}/>}></Route>
                </Router>
            );
        }
    }
}

export default App;
