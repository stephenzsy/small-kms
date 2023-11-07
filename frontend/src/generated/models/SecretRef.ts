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
/**
 * 
 * @export
 * @interface SecretRef
 */
export interface SecretRef {
    /**
     * 
     * @type {string}
     * @memberof SecretRef
     */
    id: string;
    /**
     * 
     * @type {Date}
     * @memberof SecretRef
     */
    updated: Date;
    /**
     * 
     * @type {Date}
     * @memberof SecretRef
     */
    deleted?: Date;
    /**
     * 
     * @type {string}
     * @memberof SecretRef
     */
    updatedBy?: string;
    /**
     * 
     * @type {string}
     * @memberof SecretRef
     */
    version: string;
}

/**
 * Check if a given object implements the SecretRef interface.
 */
export function instanceOfSecretRef(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "version" in value;

    return isInstance;
}

export function SecretRefFromJSON(json: any): SecretRef {
    return SecretRefFromJSONTyped(json, false);
}

export function SecretRefFromJSONTyped(json: any, ignoreDiscriminator: boolean): SecretRef {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'updated': (new Date(json['updated'])),
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'updatedBy': !exists(json, 'updatedBy') ? undefined : json['updatedBy'],
        'version': json['version'],
    };
}

export function SecretRefToJSON(value?: SecretRef | null): any {
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
        'version': value.version,
    };
}

