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
 * @interface CertificateAttributes
 */
export interface CertificateAttributes {
    /**
     * 
     * @type {string}
     * @memberof CertificateAttributes
     */
    issuer?: string;
    /**
     * 
     * @type {number}
     * @memberof CertificateAttributes
     */
    iat?: number;
    /**
     * 
     * @type {number}
     * @memberof CertificateAttributes
     */
    nbf?: number;
    /**
     * 
     * @type {number}
     * @memberof CertificateAttributes
     */
    exp?: number;
}

/**
 * Check if a given object implements the CertificateAttributes interface.
 */
export function instanceOfCertificateAttributes(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function CertificateAttributesFromJSON(json: any): CertificateAttributes {
    return CertificateAttributesFromJSONTyped(json, false);
}

export function CertificateAttributesFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateAttributes {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'issuer': !exists(json, 'issuer') ? undefined : json['issuer'],
        'iat': !exists(json, 'iat') ? undefined : json['iat'],
        'nbf': !exists(json, 'nbf') ? undefined : json['nbf'],
        'exp': !exists(json, 'exp') ? undefined : json['exp'],
    };
}

export function CertificateAttributesToJSON(value?: CertificateAttributes | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'issuer': value.issuer,
        'iat': value.iat,
        'nbf': value.nbf,
        'exp': value.exp,
    };
}

