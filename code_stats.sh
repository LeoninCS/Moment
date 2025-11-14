1#!/bin/bash

# ==============================================================================
#
# code_stats.sh - Ultimate Code Statistics Tool
#
# Version: v3.1 (Keep Progress Bar)
# Description: This version modifies the script to keep the final 100% progress
#              bar on screen after the analysis is complete, providing better
#              user feedback. It merges the logic of v2.10 with the UI of v2.3.
#
# ==============================================================================

# --- 配置与常量 ---
VERSION="v3.1 (Keep Progress Bar)"
SCRIPT_NAME=$(basename "$0")

# 颜色定义 (采用 v2.3 的样式)
RED=$'\033[0;31m'
GREEN=$'\033[0;32m'
YELLOW=$'\033[1;33m'
BLUE=$'\033[0;34m'
CYAN=$'\033[0;36m'
PURPLE=$'\033[0;35m'
NC=$'\033[0m'
BRIGHT_GREEN=$'\033[92m'
BRIGHT_YELLOW=$'\033[93m'
BRIGHT_RED=$'\033[91m'
BRIGHT_BLUE=$'\033[94m'
BRIGHT_CYAN=$'\033[96m'
BRIGHT_PURPLE=$'\033[95m'
TITLE_MAGENTA=$'\033[95m'

# 如果输出不是终端，则禁用颜色
if ! [[ -t 1 ]]; then
    RED="" GREEN="" YELLOW="" BLUE="" CYAN="" PURPLE="" NC=""
    BRIGHT_GREEN="" BRIGHT_YELLOW="" BRIGHT_RED="" BRIGHT_BLUE=""
    BRIGHT_CYAN="" BRIGHT_PURPLE="" TITLE_MAGENTA=""
fi

# --- 函数 ---

# 打印ASCII艺术标题 (采用 v2.3 的样式)
print_header() {
    echo -e "${PURPLE}"
cat << "EOF"
  ______          __          _____ __        __      
 / ____/___  ____/ /___       / ___// /_____ _/ /______
/ /   / __ \/ __  / _ \      \__ \/ __/ __ `/ __/ ___/
/ /___/ /_/ / /_/ /  __/     ___/ / /_/ /_/ / /_(__  ) 
\____/\____/\__,_/\___/     /____/\__/\__,_/\__/____/  
EOF
    echo -e "${NC}"
    echo -e "${BLUE}       Ultimate Code Statistics Tool ${VERSION}${NC}"
    echo ""
}

# 显示帮助信息 (采用 v2.3 的样式, 内容来自 v2.10)
show_help() {
    print_header
    echo -e "${GREEN}高级代码统计工具 ${VERSION}${NC}"
    echo
    echo -e "用法: ./${SCRIPT_NAME} [选项] ${YELLOW}[目录或文件]${NC}"
    echo
    echo "参数:"
    echo "  [目录或文件]    要分析的目标路径。如果未提供，则默认为当前目录 ('.')。"
    echo
    echo "选项:"
    echo "  -d, --debug     显示详细的调试信息，包括每个文件的分析过程。"
    echo "  -h, --help      显示此帮助信息并退出。"
    echo
    echo "示例:"
    echo "  ./${SCRIPT_NAME} .                # 分析当前目录"
    echo "  ./${SCRIPT_NAME} /path/to/project -d # 分析指定项目并开启调试模式"
    echo "  ./${SCRIPT_NAME} main.go            # 分析单个文件"
}

# 打印调试日志 (来自 v2.10)
debug_log() {
    if [ "$DEBUG" = "true" ]; then
        # 调试信息在进度条后换行打印，避免覆盖
        printf "\r\033[K" 
        echo -e "${YELLOW}[DEBUG]${NC} $@"
    fi
}

# 显示进度条 (来自 v2.3)
show_progress() {
    # 调试模式下不显示进度条
    if [ "$DEBUG" = "true" ]; then return; fi
    local current=$1; local total=$2; local width=50
    if [ "$total" -eq 0 ]; then return; fi
    local percentage=$((current * 100 / total))
    local completed=$((current * width / total))
    local progress_bar; progress_bar=$(printf "%${completed}s" | tr ' ' '=')
    local remaining_bar; remaining_bar=$(printf "%$((width-completed))s")
    printf "\r${CYAN}分析中: [${progress_bar}>${remaining_bar}] ${percentage}%% (${current}/${total})${NC}"
}

# 根据文件扩展名确定语言和注释风格 (逻辑来自 v2.10)
get_language_info() {
    local filename
    filename=$(basename "$1")
    local extension="${filename##*.}"

    if [[ "$filename" == "$extension" ]]; then
        extension="$filename"
    fi

    case "${extension,,}" in
        go)         echo "Go //";;
        md)         echo "Markdown";;
        sh|bash|zsh)echo "Shell #";;
        yml|yaml)   echo "YAML #";;
        py)         echo "Python #";;
        c|h)        echo "C //";;
        cpp|hpp|cxx)echo "C++ //";;
        java)       echo "Java //";;
        js|mjs|jsx) echo "JavaScript //";;
        ts|tsx)     echo "TypeScript //";;
        html|htm)   echo "HTML";;
        css)        echo "CSS";;
        rb)         echo "Ruby #";;
        rs)         echo "Rust //";;
        php)        echo "PHP //";;
        toml)       echo "TOML #";;
        Dockerfile) echo "Dockerfile #";;
        Makefile)   echo "Makefile #";;
        *)          echo "Unknown";;
    esac
}

# --- 主逻辑 ---

main() {
    # 1. 参数解析 (逻辑来自 v2.10)
    TARGET_PATH="."
    DEBUG=false

    for arg in "$@"; do
        case "$arg" in
            -h|--help)
                show_help
                exit 0
                ;;
            -d|--debug)
                DEBUG=true
                ;;
            *)
                if [[ ! "$arg" =~ ^- ]]; then
                    TARGET_PATH="$arg"
                else
                    echo -e "${RED}错误: 未知选项 '$arg'${NC}"
                    show_help
                    exit 1
                fi
                ;;
        esac
    done
    
    if [[ "$#" -eq 0 || ( "$#" -eq 1 && "$DEBUG" = true ) ]]; then
        local non_flag_args=0
        for arg in "$@"; do
            if [[ ! "$arg" =~ ^- ]]; then
                ((non_flag_args++))
            fi
        done
        if [ $non_flag_args -eq 0 ]; then
            show_help
            exit 0
        fi
    fi

    if [ ! -e "$TARGET_PATH" ]; then
        echo -e "${RED}错误: 路径 '$TARGET_PATH' 不存在。${NC}"
        exit 1
    fi
    TARGET_PATH_ABS=$(realpath "$TARGET_PATH")

    # 2. 初始化
    print_header

    echo -e "${PURPLE}=== 高级代码统计分析 ===${NC}"
    echo -e "${BLUE}分析目标: ${GREEN}$TARGET_PATH_ABS${NC}"
    if [ "$DEBUG" = true ]; then
        echo -e "${YELLOW}调试模式: 已启用${NC}"
    fi
    echo ""

    declare -A lang_files lang_total lang_blank lang_comment lang_code
    total_files=0 total_lines=0 total_blank=0 total_comment=0 total_code=0
    
    # 3. 文件扫描
    echo -e "${YELLOW}正在扫描文件...${NC}"
    local file_list
    if [ -d "$TARGET_PATH" ]; then
        readarray -t file_list < <(find "$TARGET_PATH" -type f \
            -not \( \
                -path '*/.git/*' -o -path '*/.idea/*' -o -path '*/.vscode/*' -o \
                -path '*/node_modules/*' -o -path '*/vendor/*' -o -path '*/dist/*' -o \
                -path '*/build/*' -o -path '*/target/*' -o \
                -iname '*.log' -o -iname '*.lock' -o -iname '*.min.*' -o \
                -iname '*.svg' -o -iname '*.png' -o -iname '*.jpg' -o \
                -iname '*.so' -o -iname '*.dll' -o -iname '*.exe' -o -iname '*.bin' -o -iname '*.swp' \
            \) -print)
    elif [ -f "$TARGET_PATH" ]; then
        file_list=("$TARGET_PATH")
    else
        echo -e "${RED}错误: '$TARGET_PATH' 不是一个有效的文件或目录。${NC}"
        exit 1
    fi

    local num_files_found=${#file_list[@]}
    if [ "$num_files_found" -eq 0 ]; then
        echo -e "${RED}扫描完成，未找到可分析的文件。${NC}"
        exit 0
    fi
    echo -e "${GREEN}扫描完成，找到 ${num_files_found} 个文件准备分析。${NC}"
    echo ""
    # 4. 分析每个文件
    local progress_counter=0
    local analyzed_files_count=0
    for file in "${file_list[@]}"; do
        ((progress_counter++))
        show_progress "$progress_counter" "$num_files_found"

        lang_info=($(get_language_info "$file"))
        lang=${lang_info[0]}
        comment_char=${lang_info[1]}

        if [ "$lang" == "Unknown" ]; then
            debug_log "跳过未知文件类型: $file"
            continue
        fi

        ((analyzed_files_count++))
        debug_log "分析: [$lang] $file"

        local t_lines=0 b_lines=0 c_lines=0
        if [ -n "$comment_char" ]; then
            read -r t_lines b_lines c_lines < <(awk -v C="${comment_char}" '
                BEGIN { blank=0; comment=0 }
                /^\s*$/ { blank++; next }
                { sub(/^[ \t]+/, ""); if (substr($0, 1, length(C)) == C) comment++ }
                END { print NR, blank, comment }' "$file" 2>/dev/null)
        else
            read -r t_lines b_lines c_lines < <(awk '
                BEGIN { blank=0 }
                /^\s*$/ { blank++ }
                END { print NR, blank, 0 }' "$file" 2>/dev/null)
        fi

        local co_lines=$(( t_lines - b_lines - c_lines ))
        [ $co_lines -lt 0 ] && co_lines=0

        debug_log "  -> 总行: $t_lines, 空行: $b_lines, 注释: $c_lines, 代码: $co_lines"

        ((lang_files[$lang]++))
        ((lang_total[$lang]+=t_lines))
        ((lang_blank[$lang]+=b_lines))
        ((lang_comment[$lang]+=c_lines))
        ((lang_code[$lang]+=co_lines))
    done
    
    # <<< MODIFICATION: 显示最终的100%进度条并换行, 而不是清除它
    if [ "$DEBUG" != "true" ] && [ "$num_files_found" -gt 0 ]; then
        show_progress "$num_files_found" "$num_files_found"
        echo
    fi

    # 5. 输出报告
    echo ""
    echo -e "${GREEN}--- 按语言统计 ---${NC}"
    printf "%-15s %17s %13s %11s %12s %11s\n" "语言" "文件数" "总行数" "空行" "注释" "代码"
    printf "%s\n" "--------------------------------------------------------------------------"

    local sorted_langs=""
    if [ ${#lang_files[@]} -gt 0 ]; then
        sorted_langs=$(for lang in "${!lang_files[@]}"; do echo "$lang"; done | sort)
    fi

    for lang in $sorted_langs; do
        printf "%-15s %10d %10d %10d %10d %10d\n" \
            "$lang" \
            "${lang_files[$lang]}" \
            "${lang_total[$lang]}" \
            "${lang_blank[$lang]}" \
            "${lang_comment[$lang]}" \
            "${lang_code[$lang]}"
        
        ((total_files+=lang_files[$lang]))
        ((total_lines+=lang_total[$lang]))
        ((total_blank+=lang_blank[$lang]))
        ((total_comment+=lang_comment[$lang]))
        ((total_code+=lang_code[$lang]))
    done
    
    printf "%s\n" "--------------------------------------------------------------------------"
    printf "${GREEN}%-15s %12d %10d %10d %10d %10d${NC}\n" \
        "总计" \
        "$total_files" \
        "$total_lines" \
        "$total_blank" \
        "$total_comment" \
        "$total_code"
    
    # 项目摘要信息
    if [ "$total_lines" -gt 0 ] && [ "$total_files" -gt 0 ]; then
        echo ""
        echo -e "${TITLE_MAGENTA}--- 项目摘要信息 ---${NC}"
        
        code_density=$(( total_code * 100 / total_lines ))
        color=$BRIGHT_RED; { [ "$code_density" -ge 70 ] && color=$BRIGHT_GREEN; } || { [ "$code_density" -ge 40 ] && color=$BRIGHT_YELLOW; }
        printf "${BRIGHT_BLUE}代码密度 (代码行/总行数):    %s%d%%%s\n" "$color" "$code_density" "$NC"

        if [ $(( total_code + total_comment )) -gt 0 ]; then
            comment_ratio=$(( total_comment * 100 / (total_code + total_comment) ))
            ccolor=$BRIGHT_RED; { [ "$comment_ratio" -ge 30 ] && ccolor=$BRIGHT_GREEN; } || { [ "$comment_ratio" -ge 15 ] && ccolor=$BRIGHT_YELLOW; }
            printf "${BRIGHT_CYAN}注释率 (注释行/(代码+注释)): %s%d%%%s\n" "$ccolor" "$comment_ratio" "$NC"
        fi
        
        avg_lines=$(( total_code / total_files ))
        acolor=$BRIGHT_GREEN; { [ "$avg_lines" -ge 500 ] && acolor=$BRIGHT_RED; } || { [ "$avg_lines" -ge 200 ] && acolor=$BRIGHT_YELLOW; }
        printf "${BRIGHT_GREEN}平均文件代码行数:            %s%d 行%s\n" "$acolor" "$avg_lines" "$NC"
    fi

    echo ""
    echo -e "${PURPLE}=== 代码统计分析完成 ===${NC}"
}

# 脚本入口
main "$@"

# 赋予权限: chmod +x code_stats.sh
# 使用示例: 
            #用法: ./code_stats.sh [选项] [目录或文件]

            #参数:
            #  [目录或文件]    要分析的目标路径。如果未提供，则默认为当前目录 ('.')。

            #选项:
            #  -d, --debug     显示详细的调试信息，包括每个文件的分析过程。
            #  -h, --help      显示此帮助信息并退出。

            #示例:
            #  ./code_stats.sh .                # 分析当前目录
            #  ./code_stats.sh /path/to/project -d # 分析指定项目并开启调试模式
            #  ./code_stats.sh main.go            # 分析单个文件