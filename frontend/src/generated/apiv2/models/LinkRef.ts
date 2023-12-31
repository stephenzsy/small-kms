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
 * @interface LinkRef
 */
export interface LinkRef {
    /**
     * 
     * @type {string}
     * @memberof LinkRef
     */
    id: string;
    /**
     * 
     * @type {Date}
     * @memberof LinkRef
     */
    updated: Date;
    /**
     * 
     * @type {string}
     * @memberof LinkRef
     */
    updatedBy: string;
    /**
     * 
     * @type {Date}
     * @memberof LinkRef
     */
    deleted?: Date;
    /**
     * 
     * @type {string}
     * @memberof LinkRef
     */
    displayName?: string;
    /**
     * 
     * @type {string}
     * @memberof LinkRef
     */
    linkTo: string;
}

/**
 * Check if a given object implements the LinkRef interface.
 */
export function instanceOfLinkRef(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "updatedBy" in value;
    isInstance = isInstance && "linkTo" in value;

    return isInstance;
}

export function LinkRefFromJSON(json: any): LinkRef {
    return LinkRefFromJSONTyped(json, false);
}

export function LinkRefFromJSONTyped(json: any, ignoreDiscriminator: boolean): LinkRef {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'updated': (new Date(json['updated'])),
        'updatedBy': json['updatedBy'],
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'displayName': !exists(json, 'displayName') ? undefined : json['displayName'],
        'linkTo': json['linkTo'],
    };
}

export function LinkRefToJSON(value?: LinkRef | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'updated': (value.updated.toISOString()),
        'updatedBy': value.updatedBy,
        'deleted': value.deleted === undefined ? undefined : (value.deleted.toISOString()),
        'displayName': value.displayName,
        'linkTo': value.linkTo,
    };
}

