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
import type { CertificateUsage } from './CertificateUsage';
import {
    CertificateUsageFromJSON,
    CertificateUsageFromJSONTyped,
    CertificateUsageToJSON,
} from './CertificateUsage';
import type { JwkProperties } from './JwkProperties';
import {
    JwkPropertiesFromJSON,
    JwkPropertiesFromJSONTyped,
    JwkPropertiesToJSON,
} from './JwkProperties';
import type { SubjectAlternativeNames } from './SubjectAlternativeNames';
import {
    SubjectAlternativeNamesFromJSON,
    SubjectAlternativeNamesFromJSONTyped,
    SubjectAlternativeNamesToJSON,
} from './SubjectAlternativeNames';

/**
 * 
 * @export
 * @interface CertificateInfoFields
 */
export interface CertificateInfoFields {
    /**
     * Expiration date of the certificate
     * @type {Date}
     * @memberof CertificateInfoFields
     */
    notBefore: Date;
    /**
     * 
     * @type {SubjectAlternativeNames}
     * @memberof CertificateInfoFields
     */
    subjectAlternativeNames?: SubjectAlternativeNames;
    /**
     * 
     * @type {Array<CertificateUsage>}
     * @memberof CertificateInfoFields
     */
    usages: Array<CertificateUsage>;
    /**
     * 
     * @type {string}
     * @memberof CertificateInfoFields
     */
    issuer: string;
    /**
     * 
     * @type {JwkProperties}
     * @memberof CertificateInfoFields
     */
    jwk: JwkProperties;
    /**
     * 
     * @type {string}
     * @memberof CertificateInfoFields
     */
    pem?: string;
}

/**
 * Check if a given object implements the CertificateInfoFields interface.
 */
export function instanceOfCertificateInfoFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "notBefore" in value;
    isInstance = isInstance && "usages" in value;
    isInstance = isInstance && "issuer" in value;
    isInstance = isInstance && "jwk" in value;

    return isInstance;
}

export function CertificateInfoFieldsFromJSON(json: any): CertificateInfoFields {
    return CertificateInfoFieldsFromJSONTyped(json, false);
}

export function CertificateInfoFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateInfoFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'notBefore': (new Date(json['notBefore'])),
        'subjectAlternativeNames': !exists(json, 'subjectAlternativeNames') ? undefined : SubjectAlternativeNamesFromJSON(json['subjectAlternativeNames']),
        'usages': ((json['usages'] as Array<any>).map(CertificateUsageFromJSON)),
        'issuer': json['issuer'],
        'jwk': JwkPropertiesFromJSON(json['jwk']),
        'pem': !exists(json, 'pem') ? undefined : json['pem'],
    };
}

export function CertificateInfoFieldsToJSON(value?: CertificateInfoFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'notBefore': (value.notBefore.toISOString()),
        'subjectAlternativeNames': SubjectAlternativeNamesToJSON(value.subjectAlternativeNames),
        'usages': ((value.usages as Array<any>).map(CertificateUsageToJSON)),
        'issuer': value.issuer,
        'jwk': JwkPropertiesToJSON(value.jwk),
        'pem': value.pem,
    };
}
