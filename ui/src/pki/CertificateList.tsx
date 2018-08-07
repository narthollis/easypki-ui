import React from 'react';
import { LightWeightCertificate } from 'src/api';
import CertificateListNode from 'src/pki/CertificateListNode';

interface Props {
    root: LightWeightCertificate;
}

export default class CertificateList  extends React.PureComponent<Props> {
    public render(): React.ReactNode {
        const { root } = this.props;

        return (
            <div className="my-3 p-3 bg-white rounded shadow-sm">
                <h6 className="border-bottom border-gray pb-2 mb-0">{root.name}</h6>
                {(root.children || []).map(node => (node.children != null ? (
                    <CertificateList key={node.name} root={node} />
                ) : (
                    <CertificateListNode key={node.name} node={node} />
                )))}
            </div>
        )
    }
}
