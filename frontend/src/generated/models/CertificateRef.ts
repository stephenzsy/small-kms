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
import type { NamespaceKind } from './NamespaceKind';
import {
    NamespaceKindFromJSON,
    NamespaceKindFromJSONTyped,
    NamespaceKindToJSON,
} from './NamespaceKind';
import type { ResourceKind } from './ResourceKind';
import {
    ResourceKindFromJSON,
    ResourceKindFromJSONTyped,
    ResourceKindToJSON,
} from './ResourceKind';

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
     * @type {NamespaceKind}
     * @memberof CertificateRef
     */
    namespaceKind: NamespaceKind;
    /**
     * 
     * @type {string}
     * @memberof CertificateRef
     */
    namespaceIdentifier: string;
    /**
     * 
     * @type {ResourceKind}
     * @memberof CertificateRef
     */
    resourceKind: ResourceKind;
    /**
     * 
     * @type {string}
     * @memberof CertificateRef
     */
    resourceIdentifier: string;
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
    updatedBy: string;
    /**
     * 
     * @type {string}
     * @memberof CertificateRef
     */
    x5t: string;
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
    isInstance = isInstance && "namespaceKind" in value;
    isInstance = isInstance && "namespaceIdentifier" in value;
    isInstance = isInstance && "resourceKind" in value;
    isInstance = isInstance && "resourceIdentifier" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "updatedBy" in value;
    isInstance = isInstance && "x5t" in value;
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
        
        'id': json['_id'],
        'namespaceKind': NamespaceKindFromJSON(json['namespaceKind']),
        'namespaceIdentifier': json['namespaceIdentifier'],
        'resourceKind': ResourceKindFromJSON(json['resourceKind']),
        'resourceIdentifier': json['resourceIdentifier'],
        'updated': (new Date(json['updated'])),
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'updatedBy': json['updatedBy'],
        'x5t': json['x5t'],
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
        
        '_id': value.id,
        'namespaceKind': NamespaceKindToJSON(value.namespaceKind),
        'namespaceIdentifier': value.namespaceIdentifier,
        'resourceKind': ResourceKindToJSON(value.resourceKind),
        'resourceIdentifier': value.resourceIdentifier,
        'updated': (value.updated.toISOString()),
        'deleted': value.deleted === undefined ? undefined : (value.deleted.toISOString()),
        'updatedBy': value.updatedBy,
        'x5t': value.x5t,
        'attributes': CertificateAttributesToJSON(value.attributes),
    };
}

