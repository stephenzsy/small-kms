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
 * @interface ResourceReference1
 */
export interface ResourceReference1 {
    /**
     * 
     * @type {string}
     * @memberof ResourceReference1
     */
    nid: string;
    /**
     * 
     * @type {string}
     * @memberof ResourceReference1
     */
    rid: string;
    /**
     * 
     * @type {NamespaceKind}
     * @memberof ResourceReference1
     */
    namespaceKind: NamespaceKind;
    /**
     * 
     * @type {string}
     * @memberof ResourceReference1
     */
    namespaceIdentifier: string;
    /**
     * 
     * @type {ResourceKind}
     * @memberof ResourceReference1
     */
    resourceKind: ResourceKind;
    /**
     * 
     * @type {string}
     * @memberof ResourceReference1
     */
    resourceIdentifier: string;
    /**
     * 
     * @type {Date}
     * @memberof ResourceReference1
     */
    updated: Date;
    /**
     * 
     * @type {Date}
     * @memberof ResourceReference1
     */
    deleted?: Date;
    /**
     * 
     * @type {string}
     * @memberof ResourceReference1
     */
    updatedBy: string;
}

/**
 * Check if a given object implements the ResourceReference1 interface.
 */
export function instanceOfResourceReference1(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "nid" in value;
    isInstance = isInstance && "rid" in value;
    isInstance = isInstance && "namespaceKind" in value;
    isInstance = isInstance && "namespaceIdentifier" in value;
    isInstance = isInstance && "resourceKind" in value;
    isInstance = isInstance && "resourceIdentifier" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "updatedBy" in value;

    return isInstance;
}

export function ResourceReference1FromJSON(json: any): ResourceReference1 {
    return ResourceReference1FromJSONTyped(json, false);
}

export function ResourceReference1FromJSONTyped(json: any, ignoreDiscriminator: boolean): ResourceReference1 {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'nid': json['_nid'],
        'rid': json['_rid'],
        'namespaceKind': NamespaceKindFromJSON(json['namespaceKind']),
        'namespaceIdentifier': json['namespaceIdentifier'],
        'resourceKind': ResourceKindFromJSON(json['resourceKind']),
        'resourceIdentifier': json['resourceIdentifier'],
        'updated': (new Date(json['updated'])),
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'updatedBy': json['updatedBy'],
    };
}

export function ResourceReference1ToJSON(value?: ResourceReference1 | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        '_nid': value.nid,
        '_rid': value.rid,
        'namespaceKind': NamespaceKindToJSON(value.namespaceKind),
        'namespaceIdentifier': value.namespaceIdentifier,
        'resourceKind': ResourceKindToJSON(value.resourceKind),
        'resourceIdentifier': value.resourceIdentifier,
        'updated': (value.updated.toISOString()),
        'deleted': value.deleted === undefined ? undefined : (value.deleted.toISOString()),
        'updatedBy': value.updatedBy,
    };
}

