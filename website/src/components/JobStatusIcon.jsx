import React from 'react';

function JobStatusIcon(props) {
    return (
        <a key={props.cronJob.namespace+"-"+props.job.name} href={"/namespaces/"+props.cronJob.namespace+"/cronJobs/"+props.cronJob.name+"/jobs/"+props.job.name} data-toggle="tooltip" data-placement="bottom" data-original-title={ props.job.creationTime }>
            <svg viewBox="0 0 32 32" width="20" height="20" fill="none" strokeLinecap="round" strokeLinejoin="round" strokeWidth="3">
            {(() => {
                switch(props.job.status) {
                    case "succeeded": return <path d="M2 20 L12 28 30 4" stroke="green"/>
                    case "failed": return <path d="M2 30 L30 2 M30 30 L2 2" stroke="red"/>
                    case "active": return <circle cx="16" cy="16" r="12" stroke="orange"/>
                    default: return <path d="M2 30 L30 2 M30 30 L2 2" stroke="black"/>
                }
            })()}
            { props.job.manual ? <path d="M19,31 l0,-10 m0,2 q3,-3 6,0 l 0,8 m0,-8 q3,-3 6,0 l 0,8" stroke="black" strokeWidth="2"/> : null}
            </svg>
        </a>
    )
}

export { JobStatusIcon }
