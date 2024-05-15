<template>
    <div id="chart_side" style="
        background-color:bisque;
        flex:2;
        height:100%;
        box-shadow: var(--el-box-shadow-light);
        border-radius: 8px;">
        <!-- 将原本放在js中的弹窗设置在vue中，这样能使用vue的模板创建更好看的弹窗，当鼠标悬停时改变弹窗的可变性、位置、函数名等参数，通过ref进行这些参数的传递 -->
        <el-card style="
                max-width: 600px;
                position:absolute;
                text-align: left;
                pointer-events: none;
                background-color: #f0f0f0;
            " :style="{
                visibility: displaybox === 1 ? 'visible' : 'hidden',
                top: boxposition.y + 'px',
                left: boxposition.x + 'px'
            }
                ">
            <p class="text_item"><b>file name:</b> {{ boxvalue.file_name }}</p>
            <p class="text_item"><b>function name:</b> {{ boxvalue.function_name }}</p>
            <p class="text_item"><b>frequency:</b> {{ boxvalue.frequency }}</p>
        </el-card>
    </div>
    <div>
        <el-dialog v-model="dialogVisible" width="70%" height="90%">
            <div class="modal-header">
                <h4 class="modal-title" id="myModalLabel">{{ dialogTitle }}</h4>
            </div>
            <div class="modal-body">
                <div class="row">
                    <div class="col-md-8" id="left_col" style="height: 600px; overflow: scroll;flex: 1"  ref="leftCol">
                        <pre><code class="language-clike line-numbers">{{ codeText }}</code></pre>
                    </div>
                    <div class="col-md-8" id="right_col" style="flex: 1">
                        <div class="svgContainer" id="right_svg">
                        </div>
                    </div>
                </div>
                <el-card style="
                        max-width: 600px;
                        position:fixed;
                        text-align: left;
                        pointer-events: none;
                        background-color: #f0f0f0;
                    " :style="{
                        visibility: displayboxin === 1 ? 'visible' : 'hidden',
                        top: boxpositionin.y + 'px',
                        left: boxpositionin.x + 'px'
                    }
                        ">
                    <p class="text_item"><b>start line:</b> {{ boxvaluein.start_line }}</p>
                    <p class="text_item"><b>end line:</b> {{ boxvaluein.end_line }}</p>
                    <p class="text_item"><b>frequency:</b> {{ boxvaluein.frequency }}</p>
                </el-card>
            </div>
        </el-dialog>
    </div>
</template>

<script>
import * as d3 from 'd3';
import { ref, inject } from 'vue';
import Prism from "prismjs";
import "../../css/prism.css";
const displaybox = ref(0);
const displayboxin = ref(0);
const dialogVisible = ref(false)
const dialogTitle = ref(null);
const codeText = ref(null);
const dataLine = ref(null);
const boxposition = ref({
    "x": 0,
    "y": 0
});
const boxpositionin = ref({
    "x": 0,
    "y": 0
});
const boxvalue = ref({
    "file_name": "",
    "function_name": "",
    "frequency": ""
});
const boxvaluein = ref({
    "start_line": "",
    "end_line": "",
    "frequency": ""
});
export default {
    data() {
        return {
            displaybox, boxvalue, boxposition,
            displayboxin, boxvaluein, boxpositionin,
            dialogVisible,
            dialogTitle,
            codeText,
            dataLine
        };
    },
    mounted() {
        this.buildCFG();
    },
    setup() {
        const detail_val = inject("detail_val");
        return {
            detail_val
        }
    },
    methods: {
        buildCFG() {
            var that = this;
            var height = document.getElementById('chart_side').clientHeight;
            var width = document.getElementById('chart_side').clientWidth;
            // 根据元素大小设置画布大小
            var root = d3.select("#chart_side")
                .append('svg')
                .attr('width', width)
                .attr('height', height);

            // 放到vue的部分中了
            // // 鼠标经过节点时的弹窗
            // var div = d3.select("body").append("div")
            //     .attr("class", "tooltip")
            //     .style("opacity", 0);
            // 整体结构
            /*
                <div id="chart_side">
                    <svg>
                        <defs>
                            <marker></marker>
                        </defs>
                        <g class="links">
                            <g>连线1</g>
                            <g>连线2</g>
                            <g>连线3</g>
                        </g>
                        <g class="nodes">
                            <g>节点1</g>
                            <g>节点2</g>
                            <g>节点3</g>    
                        </g>    
                    </svg>
                </div>
             */
            // 添加连线上的箭头
            root.append("svg:defs").selectAll("marker")
                .data(["end"])      // Different link/path types can be defined here
                .enter().append("svg:marker")    // This section adds in the arrows
                .attr("id", String)
                .attr("viewBox", "0 -5 10 10")
                .attr("refX", 50)
                .attr("refY", 0)
                .attr("markerWidth", 10)
                .attr("markerHeight", 6)
                .attr("orient", "auto")
                .append("svg:path")
                .attr("d", "M0,-5L10,0L0,5")
                .style("fill", "#666")
                .style("opacity", "0.7");

            // 所有的节点和边都存在一个g标签下
            var svg = root.append("g");
            var zoomListener = d3.zoom()
                .scaleExtent([0.1, 3])  // 设置缩放比例范围
                .on("zoom", zoom);

            root.call(zoomListener);  // 应用 zoom 行为到 root 上，这里的 root 应该是一个 SVG 或类似的容器元素

            function zoom(event) {  // 注意这里要接收一个 event 参数
                // console.log('event.transform',event.transform)
                // 应用当前的缩放状态到 svg 元素上
                svg.attr("transform", event.transform);
            }
            // 存储所有的节点和边
            var nodes = [];
            var links = [];
            // 所有节点存在一个g标签下
            var svg_link = svg.append("g")
                .attr("class", "links");
            // 所有的边存在一个g标签下
            var svg_node = svg.append("g")
                .attr("class", "nodes");

            var simulation = d3.forceSimulation()
                .force("link", d3.forceLink().id(function (d) { return d.id; }))
                .force("charge", d3.forceManyBody().strength(-15))
                .force("center", d3.forceCenter(width / 2, height / 2));

            var entry_function = null;
            var max_block_sum = 0;
            var explored_functions = null;
            var unexplored_functions = null;
            d3.json("/api/static.json").then(function (graph) {
                d3.json("/api/out/block_freq.json?_=" + Math.random()).then(function (_block_freq) {
                    d3.text("/api/out/fuzzer_stats?_=" + Math.random()).then(function (text) {
                        var map = {};
                        var temp_str = text.split("\n");
                        for (var i = 0; i < temp_str.length; i++) {
                            if (temp_str[i] && temp_str[i] != "") {
                                var temp = temp_str[i].split(":");
                                temp[0] = temp[0].replace(/^\s+|\s+$/g, "");
                                temp[1] = temp[1].replace(/^\s+|\s+$/g, "");
                                map[temp[0]] = temp[1];
                            }
                        }
                        var start_time = new Date(1000 * parseInt(map["start_time"]));
                        var last_update = new Date(1000 * parseInt(map["last_update"]));
                        // add_info("cycles done: ", map["cycles_done"], "cycles_done");
                        // add_info("total path: ", map["paths_total"], "total_path");
                        // add_info("unique crashes: ", map["unique_crashes"], "unique_crashes");
                        // add_info("map density: ", map["bitmap_cvg"], "map_density");
                        // add_info("execs done: ", map["execs_done"], "execs_done");
                        // add_info("max depth: ", map["max_depth"], "max_depth");
                        // add_info("unique hangs: ", map["unique_hangs"], "unique_hangs"); 
                        that.detail_val["run_time"] = show_diff_time(last_update - start_time);
                        that.detail_val["cycles_done"] = map["cycles_done"];
                        that.detail_val["total_path"] = map["paths_total"];
                        that.detail_val["unique_crashes"] = map["unique_crashes"];
                        that.detail_val["map_density"] = map["bitmap_cvg"];
                        that.detail_val["execs_done"] = map["execs_done"];
                        that.detail_val["max_depth"] = map["max_depth"];
                        that.detail_val["unique_hangs"] = map["unique_hangs"];

                    }).catch(function (error3) {
                        console.error("Error loading /api/out/fuzzer_stats?_=", error3);
                    });
                    var functions = graph.functions;
                    // console.log("functions",functions)
                    var block_freq = _block_freq.freq;
                    var basic_blocks = graph.basic_blocks;
                    for (var key in basic_blocks) {
                        basic_blocks[key]["id"] = key;
                    }

                    // find entry function
                    if (!entry_function) {
                        for (key in functions) {
                            if (functions[key]["name"] == "LLVMFuzzerTestOneInput") {
                                entry_function = key;
                                functions[key]["fx"] = width / 2;
                                functions[key]["fy"] = height / 2;
                                break;
                            }
                        }
                    }

                    // initial explore_functions and unexplored_functions
                    if (!explored_functions || !unexplored_functions) {
                        explored_functions = {};
                        unexplored_functions = {};
                        explored_functions[entry_function] = functions[entry_function];
                        functions[entry_function]["freq"] = block_freq[functions[entry_function]["entry_block"]];
                        functions[entry_function]["id"] = entry_function;
                        nodes.push(functions[entry_function]);
                        var max_freq = functions[entry_function]["freq"];
                        var queue = [entry_function];
                        while (queue.length != 0) {
                            var new_node = queue.shift();
                            var child;
                            // console.log("functions[new_node]",functions[new_node])
                            if (max_block_sum < functions[new_node]["block_sum"]) {
                                max_block_sum = functions[new_node]["block_sum"];
                            }
                            for (child in functions[new_node]["calls"]) {
                                // 对search函数进行了修改
                                search(queue, functions[new_node]["calls"][child]);
                            }
                            for (child in functions[new_node]["refs"]) {
                                // 对search函数进行了修改
                                search(queue, functions[new_node]["refs"][child]);
                            }
                        }

                        for (key in explored_functions) {
                            explored_functions[key]["color"] = Math.log2(explored_functions[key]["freq"]) / Math.log2(max_freq);
                        }
                    }
                    
                    var link_enter = svg_link
                        .selectAll("line")
                        .data(links)
                        .enter()
                        .append("line")
                        .attr('marker-end', 'url(#end)')
                        .attr('stroke', '#666')
                        .attr('stroke-opacity', 0.6)
                    .attr("stroke-width", function (d) { return Math.sqrt(d.value); });

                    var node_enter = svg_node
                        .selectAll("g")
                        .data(nodes)
                        .enter().append("g");

                    // 节点对象
                    var circles = node_enter.append("circle")
                        .attr("id", function (d) {
                            return "node" + d["id"];
                        })
                        .attr("r", function (d) {
                            return 3 + 12 * Math.log10(d["block_sum"] + 1) / Math.log10(max_block_sum + 1);
                        })
                        .attr("stroke", function (d) {
                            if (d["name"] == "LLVMFuzzerTestOneInput") {
                                return "#23459A";
                            }
                            else {
                                return "white";
                            }
                        })// 节点边缘颜色
                        .attr("stroke-width", "2px")// 节点边缘
                        .attr("fill", function (d) { return d3.interpolateGnBu(d["color"]); })// 节点颜色样式
                        .on("mouseover", function (event, d) {  // 注意这里的参数变化，现在事件对象作为第一个参数
                            that.displaybox = 1;
                            // box显示位置
                            that.boxposition.x = event.x + 10;
                            that.boxposition.y = event.y + 28;
                            // box显示内容
                            that.boxvalue.file_name = d["filename"];
                            that.boxvalue.function_name = d["name"];
                            that.boxvalue.frequency = d["freq"];
                        })
                        .on("mouseout", function () {
                            that.displaybox = 0;
                        })
                        .on("dblclick", function () {
                            d3.select("#right_svg").select("svg").remove();
                            var d = d3.select(this).datum();
                            that.$nextTick(() => {
                                open_detail(d);// 延时保证dom加载完成
                            });
                        })
                        .call(d3.drag()
                            .on("start", dragstarted)
                            .on("drag", dragged)
                            .on("end", dragended));
                    console.log(circles);
                    simulation
                        .nodes(nodes)
                        .on("tick", ticked);

                    simulation.force("link")
                        .links(links);

                    function ticked() {
                        link_enter
                            .attr("x1", function (d) { return d.source.x; })
                            .attr("y1", function (d) { return d.source.y; })
                            .attr("x2", function (d) { return d.target.x; })
                            .attr("y2", function (d) { return d.target.y; });

                        node_enter
                            .attr("transform", function (d) {
                                if (node_enter["name"] == "LLVMFuzzerTestOneInput") {
                                    return "translate(" + width / 2 + "," + height / 2 + ")";
                                }
                                else
                                    return "translate(" + d.x + "," + d.y + ")";
                            })
                    }
                    function dfs(node, map) {
                        if (node) {
                            map[node.data.data.id] = node;
                            if (node.children) {
                                node.children.forEach(function (x) {
                                    dfs(x, map);
                                });
                            }
                        }
                    }

                    function show_file(file_path, start_line, end_line) {
                        var rawFile = new XMLHttpRequest();
                        rawFile.open("GET", "/api/" + file_path, true);
                        rawFile.onreadystatechange = function () {
                            if (rawFile.readyState == 4 && rawFile.status == "200") {
                                that.codeText = rawFile.responseText;
                                that.$nextTick(() => {
                                    Prism.highlightAll()
                                })
                                that.dataline=(start_line + 1) + "-" + (end_line + 1);
                            }
                        }
                        rawFile.send(null);
                    }
                    // function show_basic_block(d) {
                    //     if (d["data"]["data"]["line"] && d["data"]["data"]["line_end"]) {
                    //         console.log("undo")
                    //     }
                    // }             
                    function open_detail(d) {
                        that.dialogTitle = "filename: " + d.filename + "    function: " + d.name;
                        that.dialogVisible = true;

                        that.$nextTick(() => {
                            var height = document.getElementById('right_col').clientHeight;
                            var width = document.getElementById('right_col').clientWidth;

                            var svg = d3.select("#right_svg")
                                .append("svg")
                                .attr("width", width)
                                .attr("height", height)
                                .attr("viewBox", `0 0 ${width} ${height}`)
                                .attr("preserveAspectRatio", "xMidYMid meet");

                            var svg_g = svg.append("g")
                                .attr("transform", "translate(0,10)");

                            var zoom = d3.zoom()
                                .scaleExtent([0.1, 3])  // 设置缩放比例范围
                                .on("zoom", (event) => {
                                    svg_g.attr("transform", event.transform);
                                });

                            svg.call(zoom);

                            var entry_block = basic_blocks[d["entry_block"]];
                            var first_node = { "data": entry_block, "children": [] };

                            var links = [];
                            var visited = Array(basic_blocks.length).fill(false);
                            var stack = [first_node];
                            visited[parseInt(entry_block["id"])] = true;

                            while (stack.length !== 0) {
                                var new_node = stack.pop();
                                new_node["data"]["successors"].forEach(function (x) {
                                    if (!visited[x]) {
                                        visited[x] = true;
                                        var tree_node = { "data": basic_blocks[x], "children": [] };
                                        new_node["children"].push(tree_node);
                                        stack.push(tree_node);
                                    }
                                });
                            }

                            var nodes = d3.hierarchy(first_node);
                            nodes = d3.tree().size([height, width])(nodes);
                            var map = {};
                            dfs(nodes, map);
                            visited = Array(basic_blocks.length).fill(false);
                            stack.push(entry_block);
                            visited[d["entry_block"]] = true;

                            while (stack.length !== 0) {
                                var new2_node = stack.pop();
                                new2_node["successors"].forEach(function (x) {
                                    if (!visited[x]) {
                                        visited[x] = true;
                                        stack.push(basic_blocks[x]);
                                    }
                                    var edge = { "source": map[parseInt(new2_node["id"])], "target": map[x] };
                                    links.push(edge);
                                });
                            }

                            var fix_height = 30;
                            var circle_r = 8;
                            nodes = nodes.descendants();
                            var max_freq = 0;
                            nodes.forEach(function (d) {
                                d.y = d.depth * fix_height;
                                d["freq"] = block_freq[parseInt(d["data"]["data"]["id"])];
                                d["color"] = d["freq"];
                                if (max_freq < d["color"]) {
                                    max_freq = d["color"];
                                }
                            });

                            var node = svg_g.selectAll(".node")
                                .data(nodes)
                                .enter().append("g")
                                .attr("class", function (d) {
                                    return "node" + (d.children ? " node--internal" : " node--leaf");
                                })
                                .attr("transform", function (d) {
                                    return "translate(" + d.x + "," + d.y + ")";
                                });

                            node.append("circle")
                                .attr("r", circle_r)
                                .attr("stroke", "black")
                                .attr("fill", function (d) {
                                    if (max_freq === 0) {
                                        return d3.interpolateYlOrRd(0);
                                    }
                                    return d3.interpolateYlOrRd(d["color"] / max_freq);
                                })
                                .on("mouseover", function (event) {
                                    var dd = d3.select(this).datum();
                                    that.displayboxin = 1;
                                    that.boxpositionin.x = event.x + 10;
                                    that.boxpositionin.y = event.y + 28;
                                    that.boxvaluein.start_line = dd["data"]["data"]["line"];
                                    that.boxvaluein.end_line = dd["data"]["data"]["line_end"];
                                    that.boxvaluein.frequency = dd["freq"];
                                })
                                .on("mouseout", function () {
                                    that.displayboxin = 0;
                                });

                            svg_g.attr("class", "links")
                                .selectAll("path")
                                .data(links)
                                .enter().append("path")
                                .attr('stroke', function (d) {
                                    return d.target.color ? d3.interpolateYlOrRd(d.target.color) : '#666';
                                })
                                .attr('fill', 'none')
                                .attr('stroke-opacity', 0.6)
                                .attr('stroke-width', 1.5)
                                .attr("d", d3.linkHorizontal()
                                    .x(function (d) {
                                        return d.x;
                                    })
                                    .y(function (d) {
                                        return d.y;
                                    }));

                            show_file(d["filename"], d["line"], d["line_end"]);
                        });
                    }


                    function search(queue, child) {
                        // console.log("child",child)
                        // console.log("functions[child]",functions[child])
                        if (functions[child]["entry_block"] >= 0 && !explored_functions[child] && !unexplored_functions[child] && functions[child]["filename"] && functions[child]["filename"].length > 0 && functions[child]["filename"][0] != "/") {
                            if (block_freq[functions[child]["entry_block"]] == 0) {
                                unexplored_functions[child] = functions[child];
                                functions[child]["freq"] = 0;
                                functions[child]["color"] = 0;
                            } else {
                                explored_functions[child] = functions[child];
                                functions[child]["freq"] = block_freq[functions[child]["entry_block"]];
                                functions[child]["id"] = child;
                                max_freq = Math.max(max_freq, functions[child]["freq"]);
                                queue.push(child);
                            }
                            functions[child]["id"] = child;
                            nodes.push(functions[child]);
                            var edge = {};
                            edge["source"] = new_node;
                            edge["target"] = child;
                            links.push(edge);
                        }
                    }

                }).catch(function (error2) {
                    console.error("Error loading /api/out/block_freq.json?_=", error2);
                });
            }).catch(function (error) {
                console.error("Error loading /api/static.json:", error);
            });

            // function centerNode(source, svg, zoomListener, viewerWidth, viewerHeight) {
            //     var t = d3.zoomTransform(svg.node());
            //     var x = -source.x;
            //     var y = -source.y;
            //     x = x * t.k + viewerWidth / 2;
            //     y = y * t.k + viewerHeight / 2;
            //     svg.transition().duration(1000)
            //         .call(zoomListener.transform, d3.zoomIdentity.translate(x, y).scale(t.k));
            // }
            function dragstarted(event, d) {
                if (!event.active) simulation.alphaTarget(0.3).restart();
                d.fx = d.x;
                d.fy = d.y;
            }

            function dragged(event, d) {
                d.fx = event.x;
                d.fy = event.y;
            }

            function dragended(event, d) {
                if (!event.active) simulation.alphaTarget(0);
                d.fx = null;
                d.fy = null;
            }

            function show_diff_time(x) {
                x = Math.floor(x / 1000);
                var hour = Math.floor(x / 3600);
                var min = Math.floor((x % 3600) / 60);
                var sec = Math.floor((x % 3600) % 60);
                return hour + " hour " + min + " min " + sec + " sec";
            }
            // function add_info(text_var, text_prop, id) {
            //     var father = document.getElementById(id);
            //     while (father.hasChildNodes()) {
            //         father.removeChild(father.firstChild);
            //     }

            //     var child_var = document.createElement('span');
            //     child_var.style.color = "#8A2BE2";
            //     child_var.style.fontWeight = "bold";
            //     child_var.style.fontSize = "20px";
            //     child_var.appendChild(document.createTextNode(text_var));
            //     father.appendChild(child_var);

            //     var child_prop = document.createElement('span');
            //     child_prop.style.color = "GhostWhite";
            //     child_prop.style.fontSize = "20px";
            //     child_prop.appendChild(document.createTextNode(text_prop));
            //     father.appendChild(child_prop);

            //     father.append(child_var);
            //     father.append(child_prop);
            // }

        }
    }
};
</script>

<style scoped>
.tooltip {
    position: absolute;
    text-align: center;
    width: 200px;
    height: 100px;
    padding: 2px;
    font: 12px sans-serif;
    background: lightsteelblue;
    border: 0px;
    border-radius: 8px;
    pointer-events: none;
    opacity: 0;
}

.text_item {
    margin: 0;
    padding: 0;
    font-size: 16px;
    z-index: 9999;
}
#left_col pre {
    overflow: initial !important;
}
.row {
    display: flex;
}
.col-md-8 {
    flex: 2;
}
.col-md-4 {
    flex: 1;
}
.svgContainer {
    width: 100%;
    height: 600px;
    /* overflow-x: scroll;
    overflow-y: hidden; */
}
/* .svgContainer svg {
    width: 100%;
    z-index: 9998;
    display: block; 
} */

.svgContainer svg {
    width: 100%;
    height: 100%;
    display: block;
}
/* .modal-dialog {
    width: 80%;
    margin: 10% auto;
}

.modal-body {
    padding: 15px;
}

.modal-content {
    border-radius: 5px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.modal-header {
    padding: 15px;
    border-bottom: 1px solid #e5e5e5;
}

.modal-footer {
    padding: 15px;
    border-top: 1px solid #e5e5e5;
} */
</style>
