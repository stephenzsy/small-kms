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
import type { KeyStatus } from './KeyStatus';
import {
    KeyStatusFromJSON,
    KeyStatusFromJSONTyped,
    KeyStatusToJSON,
} from './KeyStatus';

/**
 * 
 * @export
 * @interface KeyRef
 */
export interface KeyRef {
    /**
     * 
     * @type {string}
     * @memberof KeyRef
     */
    id: string;
    /**
     * 
     * @type {Date}
     * @memberof KeyRef
     */
    updated: Date;
    /**
     * 
     * @type {string}
     * @memberof KeyRef
     */
    updatedBy: string;
    /**
     * 
     * @type {Date}
     * @memberof KeyRef
     */
    deleted?: Date;
    /**
     * 
     * @type {string}
     * @memberof KeyRef
     */
    displayName?: string;
    /**
     * 
     * @type {KeyStatus}
     * @memberof KeyRef
     */
    status: KeyStatus;
    /**
     * 
     * @type {number}
     * @memberof KeyRef
     */
    iat: number;
    /**
     * 
     * @type {number}
     * @memberof KeyRef
     */
    exp?: number;
    /**
     * 
     * @type {string}
     * @memberof KeyRef
     */
    policyIdentifier: string;
}

/**
 * Check if a given object implements the KeyRef interface.
 */
export function instanceOfKeyRef(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "updatedBy" in value;
    isInstance = isInstance && "status" in value;
    isInstance = isInstance && "iat" in value;
    isInstance = isInstance && "policyIdentifier" in value;

    return isInstance;
}

export function KeyRefFromJSON(json: any): KeyRef {
    return KeyRefFromJSONTyped(json, false);
}

export function KeyRefFromJSONTyped(json: any, ignoreDiscriminator: boolean): KeyRef {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'updated': (new Date(json['updated'])),
        'updatedBy': json['updatedBy'],
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'displayName': !exists(json, 'displayName') ? undefined : json['displayName'],
        'status': KeyStatusFromJSON(json['status']),
        'iat': json['iat'],
        'exp': !exists(json, 'exp') ? undefined : json['exp'],
        'policyIdentifier': json['policyIdentifier'],
    };
}

export function KeyRefToJSON(value?: KeyRef | null): any {
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
        'status': KeyStatusToJSON(value.status),
        'iat': value.iat,
        'exp': value.exp,
        'policyIdentifier': value.policyIdentifier,
    };
}

