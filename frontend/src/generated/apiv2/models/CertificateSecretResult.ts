/* tslint:disable */
/* eslint-disable */
/**
 * Cryptocat API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.3
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
 * @interface CertificateSecretResult
 */
export interface CertificateSecretResult {
    /**
     * JWE encrypted certificate in PEM format
     * @type {string}
     * @memberof CertificateSecretResult
     */
    payload: string;
}

/**
 * Check if a given object implements the CertificateSecretResult interface.
 */
export function instanceOfCertificateSecretResult(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "payload" in value;

    return isInstance;
}

export function CertificateSecretResultFromJSON(json: any): CertificateSecretResult {
    return CertificateSecretResultFromJSONTyped(json, false);
}

export function CertificateSecretResultFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateSecretResult {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'payload': json['payload'],
    };
}

export function CertificateSecretResultToJSON(value?: CertificateSecretResult | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'payload': value.payload,
    };
}
