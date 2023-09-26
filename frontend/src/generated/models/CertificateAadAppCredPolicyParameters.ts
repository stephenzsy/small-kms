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
import type { CertificateIdentifier } from './CertificateIdentifier';
import {
    CertificateIdentifierFromJSON,
    CertificateIdentifierFromJSONTyped,
    CertificateIdentifierToJSON,
} from './CertificateIdentifier';

/**
 * 
 * @export
 * @interface CertificateAadAppCredPolicyParameters
 */
export interface CertificateAadAppCredPolicyParameters {
    /**
     * 
     * @type {CertificateIdentifier}
     * @memberof CertificateAadAppCredPolicyParameters
     */
    certificateIdentifier: CertificateIdentifier;
}

/**
 * Check if a given object implements the CertificateAadAppCredPolicyParameters interface.
 */
export function instanceOfCertificateAadAppCredPolicyParameters(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "certificateIdentifier" in value;

    return isInstance;
}

export function CertificateAadAppCredPolicyParametersFromJSON(json: any): CertificateAadAppCredPolicyParameters {
    return CertificateAadAppCredPolicyParametersFromJSONTyped(json, false);
}

export function CertificateAadAppCredPolicyParametersFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateAadAppCredPolicyParameters {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'certificateIdentifier': CertificateIdentifierFromJSON(json['certificateIdentifier']),
    };
}

export function CertificateAadAppCredPolicyParametersToJSON(value?: CertificateAadAppCredPolicyParameters | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'certificateIdentifier': CertificateIdentifierToJSON(value.certificateIdentifier),
    };
}

