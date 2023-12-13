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
import type { JsonWebKey } from './JsonWebKey';
import {
    JsonWebKeyFromJSON,
    JsonWebKeyFromJSONTyped,
    JsonWebKeyToJSON,
} from './JsonWebKey';

/**
 * 
 * @export
 * @interface CertificateSecretRequest
 */
export interface CertificateSecretRequest {
    /**
     * 
     * @type {JsonWebKey}
     * @memberof CertificateSecretRequest
     */
    jwk: JsonWebKey;
}

/**
 * Check if a given object implements the CertificateSecretRequest interface.
 */
export function instanceOfCertificateSecretRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "jwk" in value;

    return isInstance;
}

export function CertificateSecretRequestFromJSON(json: any): CertificateSecretRequest {
    return CertificateSecretRequestFromJSONTyped(json, false);
}

export function CertificateSecretRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateSecretRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jwk': JsonWebKeyFromJSON(json['jwk']),
    };
}

export function CertificateSecretRequestToJSON(value?: CertificateSecretRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'jwk': JsonWebKeyToJSON(value.jwk),
    };
}
