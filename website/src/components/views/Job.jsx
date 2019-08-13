import React from 'react';

class Job extends React.Component {
  render() {
    console.log(this.props)
    return (
      <div id="job">
        This page will have a Job on it.
      </div>
    );
  }
}

export default Job

