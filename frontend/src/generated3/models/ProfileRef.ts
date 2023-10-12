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

/**
 * 
 * @export
 * @interface ProfileRef
 */
export interface ProfileRef {
    /**
     * 
     * @type {string}
     * @memberof ProfileRef
     */
    id: string;
    /**
     * 
     * @type {string}
     * @memberof ProfileRef
     */
    locator: string;
    /**
     * Time when the resoruce was last updated
     * @type {Date}
     * @memberof ProfileRef
     */
    updated?: Date;
    /**
     * 
     * @type {string}
     * @memberof ProfileRef
     */
    updatedBy?: string;
    /**
     * Time when the deleted was deleted
     * @type {Date}
     * @memberof ProfileRef
     */
    deleted?: Date;
    /**
     * 
     * @type {{ [key: string]: any; }}
     * @memberof ProfileRef
     */
    metadata?: { [key: string]: any; };
    /**
     * 
     * @type {NamespaceKind}
     * @memberof ProfileRef
     */
    type: NamespaceKind;
    /**
     * Display name of the resource
     * @type {string}
     * @memberof ProfileRef
     */
    displayName: string;
    /**
     * Whether the resource is managed by the application
     * @type {boolean}
     * @memberof ProfileRef
     */
    isAppManaged?: boolean;
}

/**
 * Check if a given object implements the ProfileRef interface.
 */
export function instanceOfProfileRef(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "locator" in value;
    isInstance = isInstance && "type" in value;
    isInstance = isInstance && "displayName" in value;

    return isInstance;
}

export function ProfileRefFromJSON(json: any): ProfileRef {
    return ProfileRefFromJSONTyped(json, false);
}

export function ProfileRefFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProfileRef {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'locator': json['locator'],
        'updated': !exists(json, 'updated') ? undefined : (new Date(json['updated'])),
        'updatedBy': !exists(json, 'updatedBy') ? undefined : json['updatedBy'],
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'metadata': !exists(json, 'metadata') ? undefined : json['metadata'],
        'type': NamespaceKindFromJSON(json['type']),
        'displayName': json['displayName'],
        'isAppManaged': !exists(json, 'isAppManaged') ? undefined : json['isAppManaged'],
    };
}

export function ProfileRefToJSON(value?: ProfileRef | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'locator': value.locator,
        'updated': value.updated === undefined ? undefined : (value.updated.toISOString()),
        'updatedBy': value.updatedBy,
        'deleted': value.deleted === undefined ? undefined : (value.deleted.toISOString()),
        'metadata': value.metadata,
        'type': NamespaceKindToJSON(value.type),
        'displayName': value.displayName,
        'isAppManaged': value.isAppManaged,
    };
}

