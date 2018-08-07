import React from 'react';

import { getCertificateTree, LightWeightCertificate } from 'src/api';
import CertificateList from 'src/pki/CertificateList';
import CertificateListNode from 'src/pki/CertificateListNode';

interface Props {
}

interface State {
    isFetching: boolean;
    error?: string;
    tree: LightWeightCertificate[];
}

export default class CertificateListController extends React.PureComponent<Props, State> {
    public readonly state: Readonly<State> = {
        isFetching: true,
        tree: [],
    };

    public componentDidMount(): void {
        getCertificateTree()
            .then(tree => {
                this.setState({ tree, isFetching: false });
            })
            .catch(err => {
                this.setState({ error: err.toString(), isFetching: false });
            });
    }

    public render(): React.ReactNode {
        if (this.state.isFetching) {
            return (
                <h1>Fetching...</h1>
            );
        }

        if (this.state.error != null) {
            return <p>{this.state.error}</p>
        }

        return (
            <ul>
                {this.state.tree.map(n => <CertificateList key={n.name} root={n} />)}
            </ul>
        );
    }
}
