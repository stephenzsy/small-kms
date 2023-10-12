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
 * @interface CertificateRefFields
 */
export interface CertificateRefFields {
    /**
     * X.509 certificate SHA-1 thumbprint
     * @type {string}
     * @memberof CertificateRefFields
     */
    thumbprint: string;
    /**
     * Common name
     * @type {string}
     * @memberof CertificateRefFields
     */
    subjectCommonName: string;
    /**
     * Expiration date of the certificate
     * @type {Date}
     * @memberof CertificateRefFields
     */
    notAfter: Date;
    /**
     * 
     * @type {string}
     * @memberof CertificateRefFields
     */
    template: string;
    /**
     * Whether the certificate has been issued
     * @type {boolean}
     * @memberof CertificateRefFields
     */
    isIssued: boolean;
}

/**
 * Check if a given object implements the CertificateRefFields interface.
 */
export function instanceOfCertificateRefFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "thumbprint" in value;
    isInstance = isInstance && "subjectCommonName" in value;
    isInstance = isInstance && "notAfter" in value;
    isInstance = isInstance && "template" in value;
    isInstance = isInstance && "isIssued" in value;

    return isInstance;
}

export function CertificateRefFieldsFromJSON(json: any): CertificateRefFields {
    return CertificateRefFieldsFromJSONTyped(json, false);
}

export function CertificateRefFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateRefFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'thumbprint': json['thumbprint'],
        'subjectCommonName': json['subjectCommonName'],
        'notAfter': (new Date(json['notAfter'])),
        'template': json['template'],
        'isIssued': json['isIssued'],
    };
}

export function CertificateRefFieldsToJSON(value?: CertificateRefFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'thumbprint': value.thumbprint,
        'subjectCommonName': value.subjectCommonName,
        'notAfter': (value.notAfter.toISOString()),
        'template': value.template,
        'isIssued': value.isIssued,
    };
}

