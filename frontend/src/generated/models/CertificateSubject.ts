/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.1
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
 * @interface CertificateSubject
 */
export interface CertificateSubject {
    /**
     * 
     * @type {string}
     * @memberof CertificateSubject
     */
    commonName: string;
}

/**
 * Check if a given object implements the CertificateSubject interface.
 */
export function instanceOfCertificateSubject(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "commonName" in value;

    return isInstance;
}

export function CertificateSubjectFromJSON(json: any): CertificateSubject {
    return CertificateSubjectFromJSONTyped(json, false);
}

export function CertificateSubjectFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateSubject {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'commonName': json['commonName'],
    };
}

export function CertificateSubjectToJSON(value?: CertificateSubject | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'commonName': value.commonName,
    };
}
