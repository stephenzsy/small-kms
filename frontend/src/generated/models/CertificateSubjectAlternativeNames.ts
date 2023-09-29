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
 * @interface CertificateSubjectAlternativeNames
 */
export interface CertificateSubjectAlternativeNames {
    /**
     * 
     * @type {Array<string>}
     * @memberof CertificateSubjectAlternativeNames
     */
    dnsNames?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof CertificateSubjectAlternativeNames
     */
    emails?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof CertificateSubjectAlternativeNames
     */
    ipAddrs?: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof CertificateSubjectAlternativeNames
     */
    uris?: Array<string>;
}

/**
 * Check if a given object implements the CertificateSubjectAlternativeNames interface.
 */
export function instanceOfCertificateSubjectAlternativeNames(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function CertificateSubjectAlternativeNamesFromJSON(json: any): CertificateSubjectAlternativeNames {
    return CertificateSubjectAlternativeNamesFromJSONTyped(json, false);
}

export function CertificateSubjectAlternativeNamesFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateSubjectAlternativeNames {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'dnsNames': !exists(json, 'dns_names') ? undefined : json['dns_names'],
        'emails': !exists(json, 'emails') ? undefined : json['emails'],
        'ipAddrs': !exists(json, 'ipAddrs') ? undefined : json['ipAddrs'],
        'uris': !exists(json, 'uris') ? undefined : json['uris'],
    };
}

export function CertificateSubjectAlternativeNamesToJSON(value?: CertificateSubjectAlternativeNames | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'dns_names': value.dnsNames,
        'emails': value.emails,
        'ipAddrs': value.ipAddrs,
        'uris': value.uris,
    };
}

