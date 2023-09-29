/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface CertificateFingerprints
 */
export interface CertificateFingerprints {
    /**
     * SHA256 fingerprint of the certificate
     * @type {string}
     * @memberof CertificateFingerprints
     */
    sha256: string;
    /**
     * SHA256 fingerprint of the certificate
     * @type {string}
     * @memberof CertificateFingerprints
     */
    sha1: string;
}

/**
 * Check if a given object implements the CertificateFingerprints interface.
 */
export function instanceOfCertificateFingerprints(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "sha256" in value;
    isInstance = isInstance && "sha1" in value;

    return isInstance;
}

export function CertificateFingerprintsFromJSON(json: any): CertificateFingerprints {
    return CertificateFingerprintsFromJSONTyped(json, false);
}

export function CertificateFingerprintsFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateFingerprints {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'sha256': json['sha256'],
        'sha1': json['sha1'],
    };
}

export function CertificateFingerprintsToJSON(value?: CertificateFingerprints | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'sha256': value.sha256,
        'sha1': value.sha1,
    };
}
