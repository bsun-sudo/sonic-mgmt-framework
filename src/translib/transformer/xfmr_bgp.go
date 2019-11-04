package transformer

import (
    "errors"
    "encoding/json"
    "translib/ocbinds"
    "os/exec"
    log "github.com/golang/glog"
    ygot "github.com/openconfig/ygot/ygot"
)

func getBgpRoot (inParams XfmrParams) (*ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp, string, error) {
    pathInfo := NewPathInfo(inParams.uri)
    niName := pathInfo.Var("name")
    bgpId := pathInfo.Var("identifier")
    protoName := pathInfo.Var("name#2")
    var err error

    if len(niName) == 0 {
        return nil, "", errors.New("Network-instance-name missing")
    }

    if bgpId != "BGP" {
        return nil, "", errors.New("Network-instance-name missing")
    }

    if len(protoName) == 0 {
        return nil, "", errors.New("Network-instance-name missing")
    }

	deviceObj := (*inParams.ygRoot).(*ocbinds.Device)
    netInstsObj := deviceObj.NetworkInstances

    if netInstsObj.NetworkInstance == nil {
        return nil, "", errors.New("Network-instance-name missing")
    }

    netInstObj := netInstsObj.NetworkInstance[niName]
    if netInstObj == nil {
        return nil, "", errors.New("Network-instance-name missing")
    }

    if netInstObj.Protocols == nil || len(netInstObj.Protocols.Protocol) == 0 {
        return nil, "", errors.New("Network-instance-name missing")
    }

    var protoKey ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Key
    protoKey.Identifier = ocbinds.OpenconfigPolicyTypes_INSTALL_PROTOCOL_TYPE_BGP
    protoKey.Name = protoName
    protoInstObj := netInstObj.Protocols.Protocol[protoKey]
    if protoInstObj == nil {
        return nil, "", errors.New("Network-instance-name missing")
    }
    return protoInstObj.Bgp, niName, err
}

func init () {
    XlateFuncBind("YangToDb_bgp_gbl_tbl_key_xfmr", YangToDb_bgp_gbl_tbl_key_xfmr)
    XlateFuncBind("DbToYang_bgp_gbl_tbl_key_xfmr", DbToYang_bgp_gbl_tbl_key_xfmr)
    XlateFuncBind("YangToDb_bgp_always_compare_med_enable_xfmr", YangToDb_bgp_always_compare_med_enable_xfmr)
    XlateFuncBind("DbToYang_bgp_always_compare_med_enable_xfmr", DbToYang_bgp_always_compare_med_enable_xfmr)
    XlateFuncBind("YangToDb_bgp_allow_multiple_as_xfmr", YangToDb_bgp_allow_multiple_as_xfmr)
    XlateFuncBind("DbToYang_bgp_allow_multiple_as_xfmr", DbToYang_bgp_allow_multiple_as_xfmr)
    XlateFuncBind("YangToDb_bgp_graceful_restart_status_xfmr", YangToDb_bgp_graceful_restart_status_xfmr)
    XlateFuncBind("DbToYang_bgp_graceful_restart_status_xfmr", DbToYang_bgp_graceful_restart_status_xfmr)
    XlateFuncBind("DbToYang_bgp_nbrs_nbr_state_xfmr", DbToYang_bgp_nbrs_nbr_state_xfmr)
    XlateFuncBind("YangToDb_bgp_ignore_as_path_length_xfmr", YangToDb_bgp_ignore_as_path_length_xfmr)
    XlateFuncBind("DbToYang_bgp_ignore_as_path_length_xfmr", DbToYang_bgp_ignore_as_path_length_xfmr)
    XlateFuncBind("YangToDb_bgp_external_compare_router_id_xfmr", YangToDb_bgp_external_compare_router_id_xfmr)
    XlateFuncBind("DbToYang_bgp_external_compare_router_id_xfmr", DbToYang_bgp_external_compare_router_id_xfmr)
}

var YangToDb_bgp_gbl_tbl_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
    var err error

    pathInfo := NewPathInfo(inParams.uri)
    niName := pathInfo.Var("name")
    protoName := pathInfo.Var("name#2")

    if protoName != "bgp" {
        return niName, errors.New("Invalid protocol name : " + protoName)
    }

    log.Info("URI VRF ", niName)

    return niName, err
}

var DbToYang_bgp_gbl_tbl_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
    rmap := make(map[string]interface{})
    var err error
    entry_key := inParams.key
    log.Info("DbToYang_bgp_gbl_tbl_key: ", entry_key)

    rmap["name"] = entry_key
    return rmap, err
}

var YangToDb_bgp_always_compare_med_enable_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
    res_map := make(map[string]string)

    log.Info("YangToDb_bgp_always_compare_med_enable_xfmr: ")
    if inParams.param == nil {
	    res_map["always_compare_med"] = ""
	    return res_map, nil
	}

    enabled, _ := inParams.param.(*bool)
    var enStr string
    if *enabled == true {
        enStr = "true"
    } else {
        enStr = "false"
    }
    res_map["always_compare_med"] = enStr

    return res_map, nil
}

var DbToYang_bgp_always_compare_med_enable_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
    var err error
    result := make(map[string]interface{})

    data := (*inParams.dbDataMap)[inParams.curDb]
    log.Info("DbToYang_bgp_always_compare_med_enable_xfmr", data, "inParams : ", inParams)

    pTbl := data["BGP_GLOBALS"]
    if _, ok := pTbl[inParams.key]; !ok {
        log.Info("DbToYang_bgp_always_compare_med_enable_xfmr BGP globals not found : ", inParams.key)
        return result, errors.New("BGP globals not found : " + inParams.key)
    }
    niInst := pTbl[inParams.key]
    always_compare_med_enable, ok := niInst.Field["always_compare_med"]
    if ok {
        if always_compare_med_enable == "true" {
            result["always-compare-med"] = true
        } else {
            result["always-compare-med"] = false
        }
    } else {
        log.Info("always_compare_med field not found in DB")
    }
    return result, err
}

var YangToDb_bgp_allow_multiple_as_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
    res_map := make(map[string]string)

    log.Info("YangToDb_bgp_allow_multiple_as_xfmr: ")
    if inParams.param == nil {
	    res_map["load_balance_mp_relax"] = ""
	    return res_map, nil
	}

    allow_multiple_as, _ := inParams.param.(*bool)
    var allowMultipleAsStr string
    if *allow_multiple_as == true {
        allowMultipleAsStr = "true"
    } else {
        allowMultipleAsStr = "false"
    }
    res_map["load_balance_mp_relax"] = allowMultipleAsStr

    return res_map, nil
}

var DbToYang_bgp_allow_multiple_as_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
    var err error
    result := make(map[string]interface{})

    data := (*inParams.dbDataMap)[inParams.curDb]
    log.Info("DbToYang_bgp_allow_multiple_as_xfmr", data, "inParams : ", inParams)

    pTbl := data["BGP_GLOBALS"]
    if _, ok := pTbl[inParams.key]; !ok {
        log.Info("DbToYang_bgp_allow_multiple_as_xfmr BGP globals not found : ", inParams.key)
        return result, errors.New("BGP globals not found : " + inParams.key)
    }
    niInst := pTbl[inParams.key]
    load_balance_mp_relax_val, ok := niInst.Field["load_balance_mp_relax"]
    if ok {
        if load_balance_mp_relax_val == "true" {
            result["load_balance_mp_relax"] = true
        } else {
            result["load_balance_mp_relax"] = false
        }
    } else {
        log.Info("load_balance_mp_relax field not found in DB")
    }
    return result, err
}

var YangToDb_bgp_graceful_restart_status_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
    res_map := make(map[string]string)

    log.Info("YangToDb_bgp_graceful_restart_status_xfmr: ")
    if inParams.param == nil {
	    res_map["graceful_restart_enable"] = ""
	    return res_map, nil
	}

    gr_status, _ := inParams.param.(*bool)
    var gr_statusStr string
    if *gr_status == true {
        gr_statusStr = "true"
    } else {
        gr_statusStr = "false"
    }
    res_map["graceful_restart_enable"] = gr_statusStr

    return res_map, nil
}

var DbToYang_bgp_graceful_restart_status_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
    var err error
    result := make(map[string]interface{})

    data := (*inParams.dbDataMap)[inParams.curDb]
    log.Info("DbToYang_bgp_graceful_restart_status_xfmr", data, "inParams : ", inParams)

    pTbl := data["BGP_GLOBALS"]
    if _, ok := pTbl[inParams.key]; !ok {
        log.Info("DbToYang_bgp_graceful_restart_status_xfmr BGP globals not found : ", inParams.key)
        return result, errors.New("BGP globals not found : " + inParams.key)
    }
    niInst := pTbl[inParams.key]
    gr_enable_val, ok := niInst.Field["graceful_restart_enable"]
    if ok {
        if gr_enable_val == "true" {
            result["graceful_restart_enable"] = true
        } else {
            result["graceful_restart_enable"] = false
        }
    } else {
        log.Info("graceful_restart_enable field not found in DB")
    }
    return result, err
}

func exec_vtysh_cmd (vtysh_cmd string) (map[string]interface{}, error) {
    var err error
    oper_err := errors.New("Opertational error")

    log.Infof("Going to execute vtysh cmd ==> \"%s\"", vtysh_cmd)

    cmd := exec.Command("/usr/bin/docker", "exec", "bgp", "vtysh", "-c", vtysh_cmd)
    out_stream, err := cmd.StdoutPipe()
    if err != nil {
        log.Errorf("Can't get stdout pipe: %s\n", err)
        return nil, oper_err
    }

    err = cmd.Start()
    if err != nil {
        log.Errorf("cmd.Start() failed with %s\n", err)
        return nil, oper_err
    }

    var outputJson map[string]interface{}
    err = json.NewDecoder(out_stream).Decode(&outputJson)
    if err != nil {
        log.Errorf("Not able to decode vtysh json output: %s\n", err)
        return nil, oper_err
    }

    err = cmd.Wait()
    if err != nil {
        log.Errorf("Command execution completion failed with %s\n", err)
        return nil, oper_err
    }

    log.Infof("Successfully executed vtysh-cmd ==> \"%s\"", vtysh_cmd)
    return outputJson, err
}

func fill_nbr_state_info (nbrAddr string, nbrDataValue interface{},
                          nbr_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Neighbors_Neighbor) {
    nbrState := nbr_obj.State

    nbrState.NeighborAddress = &nbrAddr

    nbrDataJson := nbrDataValue.(map[string]interface{})

    if value, ok := nbrDataJson["bgpState"] ; ok {
        switch value {
            case "Idle":
                nbrState.SessionState = ocbinds.OpenconfigBgp_Bgp_Neighbors_Neighbor_State_SessionState_IDLE
            case "Connect":
                nbrState.SessionState = ocbinds.OpenconfigBgp_Bgp_Neighbors_Neighbor_State_SessionState_CONNECT
            case "Active":
                nbrState.SessionState = ocbinds.OpenconfigBgp_Bgp_Neighbors_Neighbor_State_SessionState_ACTIVE
            case "OpenSent":
                nbrState.SessionState = ocbinds.OpenconfigBgp_Bgp_Neighbors_Neighbor_State_SessionState_OPENSENT
            case "OpenConfirm":
                nbrState.SessionState = ocbinds.OpenconfigBgp_Bgp_Neighbors_Neighbor_State_SessionState_OPENCONFIRM
            case "Established":
                nbrState.SessionState = ocbinds.OpenconfigBgp_Bgp_Neighbors_Neighbor_State_SessionState_ESTABLISHED
            default:
                log.Infof("bgp_session_state for nbrAddr:%s ==> %s", nbrAddr, value)
        }
    }

    if value, ok := nbrDataJson["localAs"] ; ok {
        _localAs := uint32(value.(float64))
        nbrState.LocalAs = &_localAs
    }

    if value, ok := nbrDataJson["remoteAs"] ; ok {
        _peerAs := uint32(value.(float64))
        nbrState.PeerAs = &_peerAs
    }

    if value, ok := nbrDataJson["bgpTimerUpEstablishedEpoch"] ; ok {
        _lastEstablished := uint64(value.(float64))
        nbrState.LastEstablished = &_lastEstablished
    }

    if value, ok := nbrDataJson["connectionsEstablished"] ; ok {
        _establishedTransitions := uint64(value.(float64))
        nbrState.EstablishedTransitions = &_establishedTransitions
    }

    if statsMap, ok := nbrDataJson["messageStats"].(map[string]interface{}) ; ok {
        var _rcvd_msgs ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Neighbors_Neighbor_State_Messages_Received
        var _sent_msgs ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Neighbors_Neighbor_State_Messages_Sent
        var fill bool

        if value, ok := statsMap["updatesRecv"] ; ok {
            _updates_rcvd := uint64(value.(float64))
            _rcvd_msgs.UPDATE = &_updates_rcvd
            fill = true
        }
        if value, ok := statsMap["notificationsRecv"] ; ok {
            _notifs_rcvd := uint64(value.(float64))
            _rcvd_msgs.NOTIFICATION = &_notifs_rcvd
            fill = true
        }
        if value, ok := statsMap["updatesSent"] ; ok {
            _updates_sent := uint64(value.(float64))
            _sent_msgs.UPDATE = &_updates_sent
            fill = true
        }
        if value, ok := statsMap["notificationsSent"] ; ok {
            _notifs_sent := uint64(value.(float64))
            _sent_msgs.NOTIFICATION = &_notifs_sent
            fill = true
        }


        if fill {
            var _msgs ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Neighbors_Neighbor_State_Messages
            _msgs.Received = &_rcvd_msgs
            _msgs.Sent = &_sent_msgs
            nbrState.Messages = &_msgs
        }
    }
}

func get_specific_nbr_state (nbrs_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Neighbors,
                             niName string, nbrAddr string) error {
    var err error

    vtysh_cmd := "show ip bgp vrf " + niName + " neighbors " + nbrAddr + " json"
    nbrMapJson, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
        log.Errorf("Failed to fetch bgp neighbors state info for niName:%s nbrAddr:%s. Err: %s\n", niName, nbrAddr, err)
        return cmd_err
    }

    var ok bool
    var nbr_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Neighbors_Neighbor
    if len(nbrs_obj.Neighbor) == 0 {
        nbr_obj, _ = nbrs_obj.NewNeighbor (nbrAddr)
    } else {
        if nbr_obj, ok = nbrs_obj.Neighbor[nbrAddr] ; !ok {
            nbr_obj, _ = nbrs_obj.NewNeighbor (nbrAddr)
        }
    }
    ygot.BuildEmptyTree(nbr_obj)

    nbrDataJson := nbrMapJson[nbrAddr].(map[string]interface{})
    fill_nbr_state_info (nbrAddr, nbrDataJson, nbr_obj)

    return err
}

func get_all_nbr_state (nbrs_obj *ocbinds.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Protocols_Protocol_Bgp_Neighbors,
                        niName string) error {
    var err error

    vtysh_cmd := "show ip bgp vrf " + niName + " neighbors " + " json"
    nbrsMapJson, cmd_err := exec_vtysh_cmd (vtysh_cmd)
    if cmd_err != nil {
        log.Errorf("Failed to fetch all bgp neighbors state info for niName:%s. Err: %s\n", niName, err)
        return cmd_err
    }

    for nbrAddr, nbrDataJson := range nbrsMapJson {
        nbr_obj, _ := nbrs_obj.NewNeighbor (nbrAddr)
        ygot.BuildEmptyTree(nbr_obj)
        fill_nbr_state_info (nbrAddr, nbrDataJson, nbr_obj)
    }

    return err
}

var YangToDb_bgp_ignore_as_path_length_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
    res_map := make(map[string]string)

    log.Info("YangToDb_bgp_ignore_as_path_length_xfmr: ")
    if inParams.param == nil {
	    res_map["ignore_as_path_length"] = ""
	    return res_map, nil
	}

    ignore_as_path_length, _ := inParams.param.(*bool)
    var ignoreAsPathLen string
    if *ignore_as_path_length == true {
        ignoreAsPathLen = "true"
    } else {
        ignoreAsPathLen = "false"
    }
    res_map["ignore_as_path_length"] = ignoreAsPathLen

    return res_map, nil
}

var DbToYang_bgp_ignore_as_path_length_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
    var err error
    result := make(map[string]interface{})

    data := (*inParams.dbDataMap)[inParams.curDb]
    log.Info("DbToYang_bgp_ignore_as_path_length_xfmr", data, "inParams : ", inParams)

    pTbl := data["BGP_GLOBALS"]
    if _, ok := pTbl[inParams.key]; !ok {
        log.Info("DbToYang_bgp_ignore_as_path_length_xfmr BGP globals not found : ", inParams.key)
        return result, errors.New("BGP globals not found : " + inParams.key)
    }
    niInst := pTbl[inParams.key]
    ignore_as_path_length_val, ok := niInst.Field["ignore_as_path_length"]
    if ok {
        if ignore_as_path_length_val == "true" {
            result["ignore_as_path_length"] = true
        } else {
            result["ignore_as_path_length"] = false
        }
    } else {
        log.Info("ignore_as_path_length field not found in DB")
    }
    return result, err
}

var YangToDb_bgp_external_compare_router_id_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
    res_map := make(map[string]string)

    log.Info("YangToDb_bgp_external_compare_router_id_xfmr: ")
    if inParams.param == nil {
	    res_map["external_compare_router_id"] = ""
	    return res_map, nil
	}

    external_compare_router_id, _ := inParams.param.(*bool)
    var externalCompareRouterIdStr string
    if *external_compare_router_id == true {
        externalCompareRouterIdStr = "true"
    } else {
        externalCompareRouterIdStr = "false"
    }
    res_map["external_compare_router_id"] = externalCompareRouterIdStr

    return res_map, nil
}

var DbToYang_bgp_external_compare_router_id_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
    var err error
    result := make(map[string]interface{})

    data := (*inParams.dbDataMap)[inParams.curDb]
    log.Info("DbToYang_bgp_external_compare_router_id_xfmr", data, "inParams : ", inParams)

    pTbl := data["BGP_GLOBALS"]
    if _, ok := pTbl[inParams.key]; !ok {
        log.Info("DbToYang_bgp_external_compare_router_id_xfmr BGP globals not found : ", inParams.key)
        return result, errors.New("BGP globals not found : " + inParams.key)
    }
    niInst := pTbl[inParams.key]
    external_compare_router_id_val, ok := niInst.Field["external_compare_router_id"]
    if ok {
        if external_compare_router_id_val == "true" {
            result["external_compare_router_id"] = true
        } else {
            result["external_compare_router_id"] = false
        }
    } else {
        log.Info("external_compare_router_id field not found in DB")
    }
    return result, err
}

var DbToYang_bgp_nbrs_nbr_state_xfmr SubTreeXfmrDbToYang = func(inParams XfmrParams) error {
    var err error
    oper_err := errors.New("Opertational error")
    cmn_log := "GET: xfmr for BGP-nbrs state"

    bgp_obj, niName, err := getBgpRoot (inParams)
    if err != nil {
        log.Errorf ("%s failed !! Error:%s", cmn_log , err);
        return err
    }

    pathInfo := NewPathInfo(inParams.uri)
    targetUriPath, err := getYangPathFromUri(pathInfo.Path)
    nbrAddr := pathInfo.Var("neighbor-address")
    log.Infof("%s : path:%s; template:%s targetUriPath:%s niName:%s nbrAddr:%s",
              cmn_log, pathInfo.Path, pathInfo.Template, targetUriPath, niName, nbrAddr)

    supportedTgtUri := "/openconfig-network-instance:network-instances/network-instance/protocols/protocol/bgp/neighbors/neighbor/state"
    if targetUriPath != supportedTgtUri {
        log.Infof ("%s : Target-URI:%s is not %s !! Returning ok !!", cmn_log, targetUriPath, supportedTgtUri);
        return err
    }

    nbrs_obj := bgp_obj.Neighbors
    if nbrs_obj == nil {
        log.Errorf("Neighbors container missing")
        return oper_err
    }

    if len(nbrAddr) != 0 {
        err = get_specific_nbr_state (nbrs_obj, niName, nbrAddr);
    } else {
        err = get_all_nbr_state (nbrs_obj, niName);
    }

    return err;
}
