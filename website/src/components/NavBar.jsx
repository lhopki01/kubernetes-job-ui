import React from 'react';

function NavBar(props) {
  return(
    <div className="container-fluid">
      <nav className="navbar navbar-expand-lg navbar-light bg-light">
        <a className="navbar-brand" href="/">Kubernetes Job UI</a>
        <div className="collapse navbar-collapse" id="navbarSupportedContent">
          <ul className="navbar-nav">
            <li className="navbar-item">
              <a href="/cronjobs">Cronjobs</a>
            </li>
          </ul>
        </div>
      </nav>
    </div>
  )
}

export { NavBar }
