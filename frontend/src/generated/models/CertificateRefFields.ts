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
import type { CertificateAttributes } from './CertificateAttributes';
import {
    CertificateAttributesFromJSON,
    CertificateAttributesFromJSONTyped,
    CertificateAttributesToJSON,
} from './CertificateAttributes';

/**
 * 
 * @export
 * @interface CertificateRefFields
 */
export interface CertificateRefFields {
    /**
     * 
     * @type {string}
     * @memberof CertificateRefFields
     */
    thumbprint: string;
    /**
     * 
     * @type {CertificateAttributes}
     * @memberof CertificateRefFields
     */
    attributes: CertificateAttributes;
}

/**
 * Check if a given object implements the CertificateRefFields interface.
 */
export function instanceOfCertificateRefFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "thumbprint" in value;
    isInstance = isInstance && "attributes" in value;

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
        'attributes': CertificateAttributesFromJSON(json['attributes']),
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
        'attributes': CertificateAttributesToJSON(value.attributes),
    };
}

