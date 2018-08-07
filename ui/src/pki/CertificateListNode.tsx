import React from 'react';

import { RectShape } from 'react-placeholder/lib/placeholders';

import { LightWeightCertificate } from 'src/api';

interface Props {
    node: LightWeightCertificate;
}

export default class CertificateListNode extends React.PureComponent<Props> {
    public render(): React.ReactNode {
        const { node } = this.props;

        return (
            <div className="media text-muted pt-3">
                <RectShape
                    className="mr-2 rounded"
                    style={{ width: "32px", height: "32px"}}
                    color="#007bff"
                />
                <p className="media-body pb-3 mb-0 small lh-125 border-bottom border-gray">
                    <strong className="d-block text-gray-dark">{node.name}</strong>
                    Donec id elit non mi porta gravida at eget metus. Fusce dapibus, tellus ac cursus commodo,
                    tortor mauris condimentum nibh, ut fermentum massa justo sit amet risus.
                </p>
            </div>
        )
    }
}
