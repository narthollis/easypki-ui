import React from "react"
import ReactDOM from "react-dom"
import Page from 'src/page/Page';
import CertificateListController from 'src/pki/CertificateListController';

ReactDOM.render(
    <Page>
        <CertificateListController />
    </Page>,
    document.getElementById("root")
);
