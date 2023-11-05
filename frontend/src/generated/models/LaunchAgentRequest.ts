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
import type { AgentMode } from './AgentMode';
import {
    AgentModeFromJSON,
    AgentModeFromJSONTyped,
    AgentModeToJSON,
} from './AgentMode';
import type { SecretMount } from './SecretMount';
import {
    SecretMountFromJSON,
    SecretMountFromJSONTyped,
    SecretMountToJSON,
} from './SecretMount';

/**
 * 
 * @export
 * @interface LaunchAgentRequest
 */
export interface LaunchAgentRequest {
    /**
     * 
     * @type {AgentMode}
     * @memberof LaunchAgentRequest
     */
    mode: AgentMode;
    /**
     * 
     * @type {string}
     * @memberof LaunchAgentRequest
     */
    containerName: string;
    /**
     * 
     * @type {string}
     * @memberof LaunchAgentRequest
     */
    imageTag: string;
    /**
     * 
     * @type {string}
     * @memberof LaunchAgentRequest
     */
    listenerAddress: string;
    /**
     * 
     * @type {Array<string>}
     * @memberof LaunchAgentRequest
     */
    exposedPortSpecs: Array<string>;
    /**
     * 
     * @type {Array<string>}
     * @memberof LaunchAgentRequest
     */
    hostBinds: Array<string>;
    /**
     * 
     * @type {Array<SecretMount>}
     * @memberof LaunchAgentRequest
     */
    secrets: Array<SecretMount>;
    /**
     * 
     * @type {string}
     * @memberof LaunchAgentRequest
     */
    networkName?: string;
}

/**
 * Check if a given object implements the LaunchAgentRequest interface.
 */
export function instanceOfLaunchAgentRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "mode" in value;
    isInstance = isInstance && "containerName" in value;
    isInstance = isInstance && "imageTag" in value;
    isInstance = isInstance && "listenerAddress" in value;
    isInstance = isInstance && "exposedPortSpecs" in value;
    isInstance = isInstance && "hostBinds" in value;
    isInstance = isInstance && "secrets" in value;

    return isInstance;
}

export function LaunchAgentRequestFromJSON(json: any): LaunchAgentRequest {
    return LaunchAgentRequestFromJSONTyped(json, false);
}

export function LaunchAgentRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): LaunchAgentRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'mode': AgentModeFromJSON(json['mode']),
        'containerName': json['containerName'],
        'imageTag': json['imageTag'],
        'listenerAddress': json['listenerAddress'],
        'exposedPortSpecs': json['exposedPortSpecs'],
        'hostBinds': json['hostBinds'],
        'secrets': ((json['secrets'] as Array<any>).map(SecretMountFromJSON)),
        'networkName': !exists(json, 'networkName') ? undefined : json['networkName'],
    };
}

export function LaunchAgentRequestToJSON(value?: LaunchAgentRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'mode': AgentModeToJSON(value.mode),
        'containerName': value.containerName,
        'imageTag': value.imageTag,
        'listenerAddress': value.listenerAddress,
        'exposedPortSpecs': value.exposedPortSpecs,
        'hostBinds': value.hostBinds,
        'secrets': ((value.secrets as Array<any>).map(SecretMountToJSON)),
        'networkName': value.networkName,
    };
}

