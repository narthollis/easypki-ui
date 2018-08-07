import { HttpError } from 'src/api/HttpError';

const BASE_URL = 'http://localhost:8081/api/';

export interface LightWeightCertificate {
    name: string;
    commonName: string;
    notAfter: string;
    notBefore: string;
    issuer: string;

    href: string;

    children?: LightWeightCertificate[];
}

export async function getCertificateTree(): Promise<LightWeightCertificate[]> {
    const resp = await fetch(BASE_URL);

    if (!resp.ok) {
        throw new HttpError(resp.status, resp.statusText);
    }

    return resp.json();
}
