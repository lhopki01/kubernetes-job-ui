import React from 'react';
import { NavBar } from '../NavBar';

class NotFound extends React.Component {
  render() {
    return(
      <React.Fragment>
      <NavBar />
      <div className="container-fluid">
        <pre>404 Not Found</pre>
      </div>
      </React.Fragment>
    )
  }
}

export { NotFound }
