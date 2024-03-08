## Controller Card Power Control

This document describes operational use cases, rules and telemetry for power control of `CONTROLLER-CARD` components in network devices using OpenConfig.```


## Possible use cases for the gNMI approach:



1. Operator suspects faulty system behavior due to faulty hardware and wants to take it offline so they can revisit the issue in a maintenance window. But while this happens, the operator also wants to ensure that the card stays disabled across reboots so that the subject card does not take a primary role after reboot.
2. It is possible that the faulty controller card is in a boot loop and the operator finds it best to disable the card to prevent undesired state. Expectations are also that the state persists even post a system reboot.
3. Operator feels it important to keep the card shutdown to prevent unexpected outcomes post the physical swap of the faulty/alerting card so the operator can online the card in a controlled environment post replacement.
4. Following an RMA, if a new issue arises or the original problem remains, operators may instruct field techs to leave the replacement card installed. This allows for replacement upon arrival or controlled troubleshooting. In these cases, the ability to control the card's operational state is advantageous.


## Proposed Rules

Idea here is that the subject controller-card can be disabled using the OC-Path [/components/component/controller-card/config/power-admin-state](https://openconfig.net/projects/models/schemadocs/yangdoc/openconfig-platform.html#components-component-controller-card-config-power-admin-state) by setting its value to POWER\_DISABLED. Following rules:



1. Only a CONTROLLER_CARD in state redundant-roleSECONDARY will honor config/power-admin-state set toPOWER_DISABLED. If the controller-card's'redundant-role' is 'PRIMARY', and it's 'config/power-admin-state' is set to 'POWER_DISABLED', the NOS must allow the configuration. However, the change must take effect only on the next reboot. Following scenarios in context
             
`Scenario#1:`

Let's say controller-card0 is PRIMARY and controller-card1 is SECONDARY and controller-card1 receives an operator driven change of config/power-admin-state = POWER_DISABLED, then controller-card0 will power-off controller-card1 immediately.  The leaf state/power-admin-state for controller-card1 must also be POWER_DISABLED. The 'state/last-poweroff-reason/trigger' should show as USER_INITIATED. The NOS may optionally update /state/last-poweroff-reason/details. state/last-poweroff-time should record the time when the card was powered-off. For example:

```
/components/component[name=controller0]/state/redundant-role, PRIMARY
/components/component[name=controller0]/controller-card/config/power-admin-state, POWER_ENABLED
/components/component[name=controller0]/controller-card/state/power-admin-state, POWER_ENABLED
/components/component[name=controller1]/state/redundant-role, SECONDARY
/components/component[name=controller1]/controller-card/config/power-admin-state, POWER_DISABLED
/components/component[name=controller1]/controller-card/state/power-admin-state, POWER_DISABLED
/components/component[name=controller1]/state/last-poweroff-reason/trigger, USER_INITIATED
/components/component[name=controller1]/state/last-poweroff-reason/details, "User initiated Shutdown"
/components/component[name=controller1]/state/last-poweroff-time, 1706590714681765937
```

`Scenario#2:`
    
If controller-card0 is PRIMARY and controller-card1 is SECONDARY and if controller-card0 is set to config/power-admin-state = POWER_DISABLED by an operator, then controller-card0 will stay powered-on until the next reboot. The leaf state/power-admin-state must show as POWER_ENABLED. The state/last-poweroff-reason/trigger should show as USER_INITIATED. Also the NOS can update the leaf /state/last-poweroff-reason/details. The leaf state/last-poweroff-time should record the time when the card powers-off post the next reboot. For example: 

When controller-card0 is PRIMARY:

```          
/components/component[name=controller0]/state/redundant-role, PRIMARY
/components/component[name=controller0]/controller-card/config/power-admin-state, POWER_DISABLED
/components/component[name=controller0]/controller-card/state/power-admin-state, POWER_ENABLED
/components/component[name=controller1]/state/redundant-role, SECONDARY
/components/component[name=controller1]/controller-card/state/power-admin-state, POWER_ENABLED
```

After controller-card0 transitions to redundant-role SECONDARY:

```
/components/component[name=controller0]/state/redundant-role, SECONDARY
/components/component[name=controller0]/controller-card/config/power-admin-state, POWER_DISABLED
/components/component[name=controller0]/controller-card/state/power-admin-state, POWER_DISABLED
/components/component[name=controller0]/state/last-poweroff-reason/trigger, USER_INITIATED
/components/component[name=controller0]/state/last-poweroff-reason/details, "User initiated Shutdown"
/components/component[name=controller0]/state/last-poweroff-time, 1706590714681765937
/components/component[name=controller1]/state/redundant-role, PRIMARY
/components/component[name=controller1]/controller-card/config/power-admin-state, POWER_ENABLED
/components/component[name=controller1]/controller-card/state/power-admin-state, POWER_ENABLED
```

2. A controller-card which is in state/redundant-role=SECONDARY and is config/power-admin-state=POWER_DISABLED should remain powered off, even after a reboot. It is possible that after a reboot, both the controller cards are powered-on. However as soon as the configuration is loaded, the system should power-off the subject controller-card.

3. A CONTROLLER_CARD in redundant-role SECONDARY and state/power-admin-state = POWER_DISABLED cannot transition to redundant-role PRIMARY. If the PRIMARY CONTROLLER_CARD goes down, the device will be offline.

4. When a controller card boots up and loads it’s configuration, power-admin-state may be set to POWER_DISABLED. The CONTROLLER_CARD should then power off and never enter into any controller cardprimary/secondary election process. This also means that an implementation shouldn't start a controller card election process until the configuration is loaded and consumed.


## Concerns and possible failure scenarios



1. If a PRIMARY card malfunctions and ends up in a bootloop, would this approach help?

    **Response:**


    In this scenario the expectations are that the implementation takes steps to initiate a controller-card switchover operation. Therefore, the standby controller card takes over the PRIMARY role and the system stabilizes allowing for gRPC connections to be established. In this situation if the operator pushes a gNMI configuration to shutdown the SECONDARY card, the PRIMARY card must be able to power-off the SECONDARY controller-card. In this scenario the implementation (depending on their architecture) may also initiate a shutdown of the faulty card from the new PRIMARY card.

2. Since the operation relies on configuration, it is possible that the failure scenario may kick in before the configuration takes effect post a reboot.

    **Response:**


	Response here is the same as "1" above.



3. Both controller-cards are functional and the the secondary controller card is attempted to be shutdown

    **Response:**


    The shutdown can be initiated by pushing the command [/components/component[controller-card#]/controller-card/config/power-admin-state = POWER\_DISABLED](https://openconfig.net/projects/models/schemadocs/yangdoc/openconfig-platform.html#components-component-controller-card-config-power-admin-state) to the box. As a result, the secondary controller card whose [/components/component/state/redundant-role](https://openconfig.net/projects/models/schemadocs/yangdoc/openconfig-platform.html#components-component-state-redundant-role) is SECONDARY


    will be shutdown and its [/components/component[controller-card#]/controller-card/state/power-admin-state](https://openconfig.net/projects/models/schemadocs/yangdoc/openconfig-platform.html#components-component-controller-card-state-power-admin-state) will change to "POWER\_DISABLED". This configuration is expected to stay sticky in the box and when the device reboots, the subject controller-card is expected to remain Powered-off.

4. Both controller-cards are functional and the the primary controller card is attempted to be shutdown

    **Response:**


    Follow Rule#1 above

5. Only one controller-card is present and the card is attempted to be shutdown.

    **Response:**


    Since the lone controller-card is PRIMARY, please follow Rule#1 above.

6. Faulty card is shutdown and then replaced with a working card but device still has the config of `/components/component[controller-card#]/controller-card/config/power-admin-state = POWER_DISABLED`

    **Response:**


    This is as per expectations so inserting a new card into the chassis can happen anytime. However, bringing it online happens only in a controlled environment using a gNMI config push operation of `/components/component[controller-card#]/controller-card/config/power-admin-state = POWER_ENABLED`


7. Both controller-cards present during the initial boot operation. However, the controller card that is configured as "POWER\_DISABLED" boots up sooner than the other card and takes over the ACTIVE role.

    **Response:**


    The controller cards should never take PRIMARY/SECONDARY roles unless the configuration is fully loaded and processed. Hence, once the configuration is loaded, the card configured to be "POWER\_DISABLED" should shutdown. The only exception is when it is the only controller card in the chassis. In which case, the controller card can still load irrespective of its power-admin-state being POWER\_DISABLED. In this situation though, the implementation must send a Syslog message of the severity "Warning" to inform the Operator about the situation.
