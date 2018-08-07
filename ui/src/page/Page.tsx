import React from 'react';

import styled, { injectGlobal } from 'styled-components';

const Main = styled.main.attrs({ className: 'container', role: 'main' })``;

injectGlobal`
html,
body {
  overflow-x: hidden; /* Prevent scroll on narrow devices */
}

body {
  padding-top: 56px;
}

@media (max-width: 991.98px) {
  .offcanvas-collapse {
    position: fixed;
    top: 56px; /* Height of navbar */
    bottom: 0;
    left: 100%;
    width: 100%;
    padding-right: 1rem;
    padding-left: 1rem;
    overflow-y: auto;
    visibility: hidden;
    background-color: #343a40;
    transition-timing-function: ease-in-out;
    transition-duration: .3s;
    transition-property: left, visibility;
  }
  .offcanvas-collapse.open {
    left: 0;
    visibility: visible;
  }
}

.nav-scroller {
  position: relative;
  z-index: 2;
  height: 2.75rem;
  overflow-y: hidden;
}

.nav-scroller .nav {
  display: -ms-flexbox;
  display: flex;
  -ms-flex-wrap: nowrap;
  flex-wrap: nowrap;
  padding-bottom: 1rem;
  margin-top: -1px;
  overflow-x: auto;
  color: rgba(255, 255, 255, .75);
  text-align: center;
  white-space: nowrap;
  -webkit-overflow-scrolling: touch;
}

.nav-underline .nav-link {
  padding-top: .75rem;
  padding-bottom: .75rem;
  font-size: .875rem;
  color: #6c757d;
}

.nav-underline .nav-link:hover {
  color: #007bff;
}

.nav-underline .active {
  font-weight: 500;
  color: #343a40;
}

.text-white-50 { color: rgba(255, 255, 255, .5); }

.bg-purple { background-color: #6f42c1; }

.lh-100 { line-height: 1; }
.lh-125 { line-height: 1.25; }
.lh-150 { line-height: 1.5; }
`;

export default class Page extends React.PureComponent {
    public render(): React.ReactNode {
        return (
            <React.Fragment>
                <nav className="navbar navbar-expand-lg fixed-top navbar-dark bg-dark">
                    <a className="navbar-brand mr-auto mr-lg-0" href="#">EasyPKI UI</a>
                    <button className="navbar-toggler p-0 border-0" type="button" data-toggle="offcanvas">
                        <span className="navbar-toggler-icon"></span>
                    </button>

                    <div className="navbar-collapse offcanvas-collapse" id="navbarsExampleDefault">
                        <ul className="navbar-nav mr-auto">
                            <li className="nav-item active">
                                <a className="nav-link" href="#">Dashboard <span
                                    className="sr-only">(current)</span></a>
                            </li>
                            <li className="nav-item">
                                <a className="nav-link" href="#">Switch account</a>
                            </li>
                        </ul>
                        <form className="form-inline my-2 my-lg-0">
                            <input className="form-control mr-sm-2" type="text" placeholder="Search" aria-label="Search" />
                            <button className="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button>
                        </form>
                    </div>
                </nav>
                <Main>
                    {this.props.children}
                </Main>
            </React.Fragment>
        )
    }
}
