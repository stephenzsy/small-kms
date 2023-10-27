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
 * @interface ManagedAppRef
 */
export interface ManagedAppRef {
    /**
     * 
     * @type {string}
     * @memberof ManagedAppRef
     */
    id: string;
    /**
     * 
     * @type {string}
     * @memberof ManagedAppRef
     */
    uid: string;
    /**
     * 
     * @type {Date}
     * @memberof ManagedAppRef
     */
    updated: Date;
    /**
     * 
     * @type {Date}
     * @memberof ManagedAppRef
     */
    deleted?: Date;
    /**
     * 
     * @type {string}
     * @memberof ManagedAppRef
     */
    updatedBy: string;
    /**
     * 
     * @type {string}
     * @memberof ManagedAppRef
     */
    displayName: string;
    /**
     * 
     * @type {string}
     * @memberof ManagedAppRef
     */
    appId: string;
    /**
     * Object ID
     * @type {string}
     * @memberof ManagedAppRef
     */
    applicationId: string;
    /**
     * 
     * @type {string}
     * @memberof ManagedAppRef
     */
    servicePrincipalId: string;
}

/**
 * Check if a given object implements the ManagedAppRef interface.
 */
export function instanceOfManagedAppRef(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "uid" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "updatedBy" in value;
    isInstance = isInstance && "displayName" in value;
    isInstance = isInstance && "appId" in value;
    isInstance = isInstance && "applicationId" in value;
    isInstance = isInstance && "servicePrincipalId" in value;

    return isInstance;
}

export function ManagedAppRefFromJSON(json: any): ManagedAppRef {
    return ManagedAppRefFromJSONTyped(json, false);
}

export function ManagedAppRefFromJSONTyped(json: any, ignoreDiscriminator: boolean): ManagedAppRef {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'uid': json['uid'],
        'updated': (new Date(json['updated'])),
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'updatedBy': json['updatedBy'],
        'displayName': json['displayName'],
        'appId': json['appId'],
        'applicationId': json['applicationId'],
        'servicePrincipalId': json['servicePrincipalId'],
    };
}

export function ManagedAppRefToJSON(value?: ManagedAppRef | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'uid': value.uid,
        'updated': (value.updated.toISOString()),
        'deleted': value.deleted === undefined ? undefined : (value.deleted.toISOString()),
        'updatedBy': value.updatedBy,
        'displayName': value.displayName,
        'appId': value.appId,
        'applicationId': value.applicationId,
        'servicePrincipalId': value.servicePrincipalId,
    };
}

