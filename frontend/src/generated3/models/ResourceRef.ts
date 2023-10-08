/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API Shared
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
 * @interface ResourceRef
 */
export interface ResourceRef {
    /**
     * Identifier of the resource
     * @type {string}
     * @memberof ResourceRef
     */
    id: string;
    /**
     * 
     * @type {string}
     * @memberof ResourceRef
     */
    locator: string;
    /**
     * Time when the resoruce was last updated
     * @type {Date}
     * @memberof ResourceRef
     */
    updated?: Date;
    /**
     * 
     * @type {string}
     * @memberof ResourceRef
     */
    updatedBy?: string;
    /**
     * Time when the deleted was deleted
     * @type {Date}
     * @memberof ResourceRef
     */
    deleted?: Date;
    /**
     * 
     * @type {{ [key: string]: string; }}
     * @memberof ResourceRef
     */
    metadata?: { [key: string]: string; };
}

/**
 * Check if a given object implements the ResourceRef interface.
 */
export function instanceOfResourceRef(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "locator" in value;

    return isInstance;
}

export function ResourceRefFromJSON(json: any): ResourceRef {
    return ResourceRefFromJSONTyped(json, false);
}

export function ResourceRefFromJSONTyped(json: any, ignoreDiscriminator: boolean): ResourceRef {
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
    };
}

export function ResourceRefToJSON(value?: ResourceRef | null): any {
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
    };
}

