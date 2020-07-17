#!/usr/bin/env bash
set -o pipefail
#set -o nounset #This extra check has been disabled to avoid BASH 4.x crashing with some array boundaries.

#GLOBAL STUFF
VERSION="2.6"
BIN="/usr/bin"
CHANDEF=(stable fast) #Default list of channels, modify only this one if needed.


#INFO = (
#          author      => 'Pedro Amoedo'
#          contact     => 'pamoedom@redhat.com'
#          name        => 'ocp4upc.sh',
#          usage       => '(see function usage below)',
#          description => 'OCP4 Upgrade Paths Checker',
#          changelog   => '(see CHANGELOG file or git log)'
#       );

#ARGs DESCRIPTION
#$1=source_release to be used as starting point
#$2=architecture (optional), default is amd64

#USAGE
function usage()
{
  ${BIN}/echo "-------------------------------------------------------------------"
  ${BIN}/echo "OCP4 Upgrade Paths Checker ($(${BIN}/echo "${CHANDEF[@]}")) v${VERSION}"
  ${BIN}/echo ""
  ${BIN}/echo "Usage:"
  ${BIN}/echo "$0 source_version [arch]"
  ${BIN}/echo ""
  ${BIN}/echo "Source Version:"
  ${BIN}/echo "4.x        Extract default graphs using same-minor channels, e.g. '4.2'"
  ${BIN}/echo "4.x.z      Generate upgrade paths using next-minor channels, e.g. '4.2.26'"
  ${BIN}/echo "4.x.z.     Generate upgrade paths using same-minor channels, e.g. '4.2.26.'"
  ${BIN}/echo ""
  ${BIN}/echo "Arch (optional):"
  ${BIN}/echo "amd64      x86_64 (default)"
  ${BIN}/echo "s390x      IBM System/390"
  ${BIN}/echo "ppc64le    POWER8 little endian"
  ${BIN}/echo "-------------------------------------------------------------------"
  exit 1
}

#VARIABLES ($filename,$args[])
function declare_vars()
{
  cmd="$1"
  args=("$@")

  ##URLs
  GPH='https://api.openshift.com/api/upgrades_info/v1/graph'
  REL='https://quay.io/api/v1/repository/openshift-release-dev/ocp-release'

  ##ARGs
  [[ ${#args[1]} -lt 3 ]] && usage || VER=${args[1]}
  [[ -z ${args[2]-} ]] && ARC="amd64" || ARC=${args[2]}

  ##Target channel calculation & mode detection
  ! [[ ${VER} =~ ^[0-9]([.][0-9]+).*$ ]] && usage
  MAJ=$(${BIN}/echo ${VER} | ${BIN}/cut -d. -f1)
  MIN=$(${BIN}/echo ${VER} | ${BIN}/cut -d. -f2)
  ERT=$(${BIN}/echo ${VER} | ${BIN}/cut -d. -f3-) #errata version provided?
  if [ "${ERT}" != "" ]; then
    [[ ${ERT} =~ ^[0-9]+$ ]] && TRG="${MAJ}.$(( ${MIN} + 1 ))" && MOD="4.x.z"
    [[ ${ERT} =~ ^[0-9]+[.]$ ]] && TRG="${MAJ}.${MIN}" && VER=$(${BIN}/echo ${VER} | ${BIN}/cut -d. -f1,2,3) && MOD="4.x.z."
    [[ ${TRG} = "" ]] && usage
  else
    VER=$(${BIN}/echo ${VER} | ${BIN}/cut -d. -f1,2)
    TRG="${VER}"
    MOD="4.x"
  fi

  ##Edge & Node colors
  EDGs="blue" #source edges -> *
  EDGt="red" #source edges -> (LTS)
  NODs="salmon" #source node
  NODt="yellowgreen" #target nodes (LTS)
  NODi="lightgrey" #indirect nodes
  DEF="grey" #default color (keep it different)

  ##Various Arrays
  CHA=("${CHANDEF[@]}") #array of channels
  REQ=(curl jq dot) #array of pre-requisities
  RES=() #array of resulting channels in case of discard
  for chan in "${CHA[@]}"; do declare -a "LTS_${chan}=()"; done #arrays of possible target releases per channel.
  IND=() #array of indirect nodes (if any)
  EXT=() #array of direct nodes (if any)

  ##Ansi colors
  OK="${BIN}/echo -en \\033[1;32m" #green
  ERROR="${BIN}/echo -en \\033[1;31m" #red
  WARN="${BIN}/echo -en \\033[1;33m" #yellow
  INFO="${BIN}/echo -en \\033[1;34m" #blue
  NORM="${BIN}/echo -en \\033[0;39m" #default

  ##Misc
  PTH="/tmp/${cmd##*/}" #generate the tmp folder based on the current script name
  RELf="ocp4-releases.json"
  KEY='  Key \[rank=sink,shape=none,margin=0\.3,label=< <TABLE BORDER="1" STYLE="DOTTED" CELLBORDER="0" CELLSPACING="1" CELLPADDING="0"><TR><TD COLSPAN="2"><B>Key<\/B><\/TD><\/TR><TR><TD align="left">Direct Path<\/TD><TD><FONT COLOR="'"${EDGt}"'">\&\#10230\;<\/FONT><\/TD><\/TR><TR><TD align="left">Indirect Path<\/TD><TD><FONT COLOR="'"${EDGs}"'">\&\#8594\; \&\#10511\;<\/FONT><\/TD><\/TR><\/TABLE> >\]; }'
}

#PRETTY PRINT ($type_of_msg,$string,$echo_opts)
function cout()
{
  [[ -z ${3-} ]] && opts="" || opts=$3
  ${BIN}/echo -n "[" && eval "\$$1" && ${BIN}/echo -n "$1" && $NORM && ${BIN}/echo -n "] " && ${BIN}/echo ${opts} "$2"
}

#PREREQUISITES
function check_prereq()
{
  ##all tools available?
  cout "INFO" "Checking prerequisites ($(${BIN}/echo "${REQ[@]}"))... " "-n"
  for tool in "${REQ[@]}"; do ${BIN}/which ${tool} &>/dev/null; [ $? -ne 0 ] && cout "ERROR" "'${tool}' not present. Aborting execution." && exit 1; done

  ##tmp folder writable?
  if [ -d ${PTH} ]; then
    ${BIN}/touch ${PTH}/test; [ $? -ne 0 ] && cout "ERROR" "Unable to write in ${PTH}. Aborting execution." && exit 1
    ${BIN}/rm ${PTH}/*.json ${PTH}/*.gv > /dev/null 2>&1
  else
    ${BIN}/mkdir ${PTH}; [ $? -ne 0 ] && cout "ERROR" "Unable to write in ${PTH}. Aborting execution." && exit 1
  fi
  cout "OK" ""
}

#RELEASE CHECKING
function check_release()
{
  ${BIN}/curl -sH 'Accept:application/json' "${REL}" | ${BIN}/jq . > ${PTH}/${RELf}
  if [ $? -ne 0 ]; then
    cout "WARN" "Unable to curl '${REL}'"
    cout "WARN" "Do you want to continue omitting sanity checks? (y/N):" "-n"
    read -t 10 yn
    [[ ${yn} =~ ^([yY][eE][sS]|[yY])$ ]] && return || cout "ERROR" "Execution interrupted, try again later." && exit 1;
  fi
  
  if [ "${ERT}" = "" ]; then
    cout "INFO" "Checking if '${VER}' (${ARC}) has valid channels... " "-n"
    if [ "${ARC}" = "amd64" ]; then
      ${BIN}/grep "\"${VER}.*-x86_64\"" ${PTH}/${RELf} &>/dev/null; [ $? -ne 0 ] && cout "ERROR" "" && exit 1
    else
      ${BIN}/grep "\"${VER}.*-${ARC}\"" ${PTH}/${RELf} &>/dev/null; [ $? -ne 0 ] && cout "ERROR" "" && exit 1
    fi
  else
    cout "INFO" "Checking if '${VER}' (${ARC}) is a valid release... " "-n"
    if [ "${ARC}" = "amd64" ]; then
      ${BIN}/grep "\"${VER}-x86_64\"" ${PTH}/${RELf} &>/dev/null;
      ##for amd64 make an extra attempt without -x86_64 because old releases don't have any suffix
      if [ $? -ne 0 ]; then
        ${BIN}/grep "\"${VER}\"" ${PTH}/${RELf} &>/dev/null; [ $? -ne 0 ] && cout "ERROR" "" && cout "INFO" "TIP: run the script without parameters to see other available architectures." && exit 1
      fi
    else
      ${BIN}/grep "\"${VER}-${ARC}\"" ${PTH}/${RELf} &>/dev/null; [ $? -ne 0 ] && cout "ERROR" "" && cout "INFO" "TIP: run the script without parameters to see other available architectures." && exit 1
    fi
  fi
  cout "OK" ""
}

#OBTAIN UPGRADE PATHS JSONs
function get_paths()
{
  i=0
  for chan in "${CHA[@]}"; do
    ${BIN}/curl -sH 'Accept:application/json' "${GPH}?channel=${chan}-${TRG}&arch=${ARC}" > ${PTH}/${chan}-${TRG}.json
    [[ $? -ne 0 ]] && cout "ERROR" "Unable to curl '${GPH}?channel=${chan}-${TRG}&arch=${ARC}'" && cout "ERROR" "Execution interrupted, try again later." && exit 1
    ##discard empty channels
    ${BIN}/echo -n '{"nodes":[],"edges":[]}' | ${BIN}/diff ${PTH}/${chan}-${TRG}.json - &>/dev/null; [ $? -eq 0 ] && cout "WARN" "Skipping channel '${chan}-${TRG}_${ARC}', it's empty." && continue
    ##discard duplicated channels
    [[ $i -ne 0 ]] && ${BIN}/diff ${PTH}/${chan}-${TRG}.json ${PTH}/${CHA[$(( $i - 1 ))]}-${TRG}.json &>/dev/null; [ $? -eq 0 ] && cout "WARN" "Discarding channel '${chan}-${TRG}_${ARC}', it doesn't differ from '${CHA[$(( $i - 1 ))]}-${TRG}_${ARC}'." && continue
    RES=("${RES[@]}" "${chan}")
    (( i++ ))
  done
  ##reset channel list accordingly or abort if none available
  CHA=("${RES[@]}")
  if [ ${#CHA[@]} -eq 0 ]; then
    ##allow the user to make a 2nd run on the same minor channels if targeting a non-released version (corner case)
    if [ "${ERT}" != "" ] && [ "${TRG}" != "${MAJ}.${MIN}" ]; then
      cout "WARN" "You are targeting void '${TRG}' channels, do you want to re-target to '"${MAJ}.${MIN}"' instead (4.x.z. mode)? (y/N):" "-n"
      read -t 10 yn
      [[ ${yn} =~ ^([yY][eE][sS]|[yY])$ ]] && return 2 
    fi
    cout "ERROR" "There are no channels to process. Aborting execution." && exit 1;
  fi
}

#CAPTURE TARGETS ##TODO: do this against the API instead?
function capture_lts()
{
  for chan in "${CHA[@]}"; do
    var="LTS_${chan}"
    eval "${var}"="("$(${BIN}/cat ${PTH}/${chan}-${TRG}.json | ${BIN}/jq . | ${BIN}/grep "\"${TRG}." | ${BIN}/cut -d'"' -f4 | ${BIN}/sort -urV | ${BIN}/xargs)")"
  done
}

#JSON to GV
function json2gv()
{
  ##prepare the raw jq filter
  JQ_SCRIPT=$(${BIN}/echo '"digraph TITLE {\n  labelloc=c;\n  rankdir=BT;\n  label=CHANNEL" as $header |
    (
      [
        .nodes |
        to_entries[] |
        "  " + (.key | tostring) +
          " [ label=\"" + .value.version + "\"" + (
            if .value.metadata.url then ",url=\"" + .value.metadata.url + "\"" else "" end
	    ) + (')
  if [ "${ERT}" != "" ]; then
    JQ_SCRIPT=${JQ_SCRIPT}$(${BIN}/echo '            if .value.version == "'"${VER}"'" then ",shape=polygon,sides=5,peripheries=2,style=filled,color='"${NODs}"'"')
  else
    JQ_SCRIPT=${JQ_SCRIPT}$(${BIN}/echo '            if .value.version == "" then ",shape=polygon,sides=5,peripheries=3,style=filled,color='"${NODs}"'"')
  fi
  JQ_SCRIPT=${JQ_SCRIPT}$(${BIN}/echo '            elif .value.version >= "'"${TRG}"'" then ",shape=square,style=filled,color='"${DEF}"'"
            else ",shape=ellipse,style=filled,color='"${DEF}"'"
            end
          ) +
          " ];"
      ] | join("\n")
    ) as $nodes |
    (
      [
        .edges[] |
        "  " + (.[0] | tostring) + "->" + (.[1] | tostring) + ";"
      ] | join("\n")
    ) as $edges |
    [$header, $nodes, $edges, "}"] | join("\n")')

  ##generate the gv files
  for chan in "${CHA[@]}"; do
    ${BIN}/jq -r "${JQ_SCRIPT}" ${PTH}/${chan}-${TRG}.json > ${PTH}/${chan}-${TRG}.gv
    [[ $? -ne 0 ]] && cout "ERROR" "Unable to create ${PTH}/${chan}-${TRG}.gv file. Aborting execution." && exit 1
  done
}

#DISCARD CHANNELS & COLORIZE EDGES ##TODO: move this logic into JQ_SCRIPT?
function colorize()
{
  RES=() #re-initialize the array in case of channel discarding (4.x.z mode)
  for chan in "${CHA[@]}"; do
    var="LTS_${chan}"
    arr=${var}[@]
    posV=$(grep "\"${VER}\"" ${PTH}/${chan}-${TRG}.gv | awk {'print $1'})
    [[ "${posV}" = "" ]] && cout "WARN" "Skipping channel '${chan}-${TRG}_${ARC}', version '${VER}' not found." && continue
    [[ -z ${!arr-} ]] && cout "WARN" "Skipping channel '${chan}-${TRG}_${ARC}', no upgrade paths available." && continue

    ##capture list of outgoing edges (possible indirect nodes)
    IND=($(grep "\s\s${posV}->" ${PTH}/${chan}-${TRG}.gv | ${BIN}/cut -d">" -f2 | ${BIN}/cut -d";" -f1))

    ##colorize EXT->LTS edges
    for target in "${!arr}"; do
      posT=$(grep "\"${target}\"" ${PTH}/${chan}-${TRG}.gv | awk {'print $1'})
      if [ "${posT}" != "" ]; then
        for node in "${IND[@]}"; do
          ##Direct edges
          if [ "${node}" = "${posT}" ]; then
            ${BIN}/sed -i -e 's/^\(\s\s'"${posV}"'->'"${posT}"'\)\;$/\1 [color='"${EDGt}"'\,style=bold];/' ${PTH}/${chan}-${TRG}.gv
            ${BIN}/sed -i -e 's/^\(\s\s'"${posT}"'\s.*\),color=.*$/\1,color='"${NODt}"' ]\;/' ${PTH}/${chan}-${TRG}.gv
            continue
          fi
          ##Indirect edges
          ###grep is needed here because sed doesn't return a different exit code when matching
          ${BIN}/grep "\s\s${node}->${posT};" ${PTH}/${chan}-${TRG}.gv &>/dev/null
          if [ $? -eq 0 ]; then
            ##if match, colorize indirect node, indirect edge & target node at the same time (triple combo)
            ${BIN}/sed -i -e 's/^\(\s\s'"${node}"'\s.*\),color=.*$/\1,color='"${NODi}"' ]\;/;s/^\(\s\s'"${node}"'->'"${posT}"'\)\;$/\1 [color='"${EDGs}"',style=dashed];/;s/^\(\s\s'"${posT}"'\s.*\),color=.*$/\1,color='"${NODt}"' ]\;/' ${PTH}/${chan}-${TRG}.gv
            ##save final list of indirect nodes to be used below for pending source edges
            EXT=("${EXT[@]}" "${node}")
          fi
        done
      fi
    done
    ##colorize rest of source edges not yet processed
    for node in "${EXT[@]}"; do ${BIN}/sed -i -e 's/^\(\s\s'"${posV}"'->'"${node}"'\)\;$/\1 [color='"${EDGs}"',style=filled];/' ${PTH}/${chan}-${TRG}.gv; done

    ##remove non involved nodes+edges to simplify the graph
    ${BIN}/sed -i -e '/color='"${DEF}"'/d;/[0-9]\;$/d' ${PTH}/${chan}-${TRG}.gv

    ##include the graph legend
    ${BIN}/sed -i -e 's/^}$/'"${KEY}"'/' ${PTH}/${chan}-${TRG}.gv

    ##save resulting channels for subsequent operations
    RES=("${RES[@]}" "${chan}")
  done

  ##abort if the provided release is not present within any of the channels
  [[ ${#RES[@]} -eq 0 ]] && cout "ERROR" "Version '${VER}' not found (or not upgradable) within '${TRG}' channels. Aborting execution." && cout "INFO" "TIP: run the script without parameters to see other available modes." && exit 1
}

#LABELING
function label()
{
  for chan in "${RES[@]}"; do ${BIN}/sed -i -e 's/TITLE/'"${chan}"'/' ${PTH}/${chan}-${TRG}.gv; done
  for chan in "${RES[@]}"; do ${BIN}/sed -i -e 's/CHANNEL/"'"${chan}"'-'"${TRG}"'_'"${ARC}"' \('"$(${BIN}/date --rfc-3339=date)"'\)"/' ${PTH}/${chan}-${TRG}.gv; done
}

#DRAW & EXPORT
function draw()
{
  if [ "${ERT}" != "" ]; then
    for chan in "${RES[@]}"; do
      ${BIN}/dot -Tsvg ${PTH}/${chan}-${TRG}.gv -o ${chan}-${TRG}_${ARC}_$(date +%Y%m%d).svg
      [[ $? -ne 0 ]] && cout "ERROR" "Unable to export the results. Aborting execution." && exit 1 || cout "INFO" "Result exported as '${chan}-${TRG}_${ARC}_$(date +%Y%m%d).svg'"
    done
  else
    for chan in "${RES[@]}"; do
      ${BIN}/dot -Tsvg ${PTH}/${chan}-${TRG}.gv -o ${chan}-${TRG}_${ARC}_def_$(date +%Y%m%d).svg 
      [[ $? -ne 0 ]] && cout "ERROR" "Unable to export the results. Aborting execution." && exit 1 || cout "INFO" "Result exported as '${chan}-${TRG}_${ARC}_def_$(date +%Y%m%d).svg'"
    done
  fi
}

#SCRIPT WORKFLOW ($args[])
function main()
{
  args=("$@")
  declare_vars "$0" "${args[@]}"
  check_prereq
  if [ "${ERT}" != "" ]; then
    cout "INFO" "Errata provided (${MOD} mode), targeting '${TRG}' channels for upgrade path generation."
    check_release
    get_paths
    #Allow a 2nd get_paths call if same minor upgrade path generation has been selected due to empty channels (corner case)
    if [ $? -eq 2 ]; then
      CHA=("${CHANDEF[@]}")
      TRG="${MAJ}.${MIN}"
      MOD="4.x.z."
      cout "INFO" "Errata provided (${MOD} mode), targeting '${TRG}' channels for upgrade path generation."
      get_paths
    fi
    capture_lts
    json2gv
    colorize
    label
    draw
  else
    cout "INFO" "No errata provided (${MOD} mode), extracting default '${TRG}' channels."
    check_release
    get_paths
    json2gv
    label
    draw
  fi
}

#STARTING POINT
[[ $# -lt 1 ]] && usage
main "$@"

#EOF
exit 0
