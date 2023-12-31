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
 * @interface CertificateRef
 */
export interface CertificateRef {
    /**
     * 
     * @type {string}
     * @memberof CertificateRef
     */
    id: string;
    /**
     * 
     * @type {Date}
     * @memberof CertificateRef
     */
    updated: Date;
    /**
     * 
     * @type {Date}
     * @memberof CertificateRef
     */
    deleted?: Date;
    /**
     * 
     * @type {string}
     * @memberof CertificateRef
     */
    updatedBy?: string;
    /**
     * 
     * @type {string}
     * @memberof CertificateRef
     */
    thumbprint: string;
    /**
     * 
     * @type {CertificateAttributes}
     * @memberof CertificateRef
     */
    attributes: CertificateAttributes;
}

/**
 * Check if a given object implements the CertificateRef interface.
 */
export function instanceOfCertificateRef(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "thumbprint" in value;
    isInstance = isInstance && "attributes" in value;

    return isInstance;
}

export function CertificateRefFromJSON(json: any): CertificateRef {
    return CertificateRefFromJSONTyped(json, false);
}

export function CertificateRefFromJSONTyped(json: any, ignoreDiscriminator: boolean): CertificateRef {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'updated': (new Date(json['updated'])),
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'updatedBy': !exists(json, 'updatedBy') ? undefined : json['updatedBy'],
        'thumbprint': json['thumbprint'],
        'attributes': CertificateAttributesFromJSON(json['attributes']),
    };
}

export function CertificateRefToJSON(value?: CertificateRef | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'updated': (value.updated.toISOString()),
        'deleted': value.deleted === undefined ? undefined : (value.deleted.toISOString()),
        'updatedBy': value.updatedBy,
        'thumbprint': value.thumbprint,
        'attributes': CertificateAttributesToJSON(value.attributes),
    };
}

