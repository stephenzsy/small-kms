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
import type { NamespaceType } from './NamespaceType';
import {
    NamespaceTypeFromJSON,
    NamespaceTypeFromJSONTyped,
    NamespaceTypeToJSON,
} from './NamespaceType';

/**
 * 
 * @export
 * @interface NamespaceRef
 */
export interface NamespaceRef {
    /**
     * Unique ID of the namespace
     * @type {string}
     * @memberof NamespaceRef
     */
    namespaceId: string;
    /**
     * 
     * @type {string}
     * @memberof NamespaceRef
     */
    id: string;
    /**
     * Unique ID of the user who created the policy
     * @type {string}
     * @memberof NamespaceRef
     */
    updatedBy: string;
    /**
     * Time when the policy was last updated
     * @type {Date}
     * @memberof NamespaceRef
     */
    updated: Date;
    /**
     * 
     * @type {NamespaceType}
     * @memberof NamespaceRef
     */
    objectType: NamespaceType;
    /**
     * 
     * @type {string}
     * @memberof NamespaceRef
     */
    displayName: string;
}

/**
 * Check if a given object implements the NamespaceRef interface.
 */
export function instanceOfNamespaceRef(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "namespaceId" in value;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updatedBy" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "objectType" in value;
    isInstance = isInstance && "displayName" in value;

    return isInstance;
}

export function NamespaceRefFromJSON(json: any): NamespaceRef {
    return NamespaceRefFromJSONTyped(json, false);
}

export function NamespaceRefFromJSONTyped(json: any, ignoreDiscriminator: boolean): NamespaceRef {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'namespaceId': json['namespaceId'],
        'id': json['id'],
        'updatedBy': json['updatedBy'],
        'updated': (new Date(json['updated'])),
        'objectType': NamespaceTypeFromJSON(json['objectType']),
        'displayName': json['displayName'],
    };
}

export function NamespaceRefToJSON(value?: NamespaceRef | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'namespaceId': value.namespaceId,
        'id': value.id,
        'updatedBy': value.updatedBy,
        'updated': (value.updated.toISOString()),
        'objectType': NamespaceTypeToJSON(value.objectType),
        'displayName': value.displayName,
    };
}

