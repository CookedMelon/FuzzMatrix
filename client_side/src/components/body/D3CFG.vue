<template>
    <div id="chart_side" style="
        background-color:bisque;
        flex:2;
        height:100%;
        box-shadow: var(--el-box-shadow-light);
        border-radius: 8px;"
    >
        <!-- 将原本放在js中的弹窗设置在vue中，这样能使用vue的模板创建更好看的弹窗，当鼠标悬停时改变弹窗的可变性、位置、函数名等参数，通过ref进行这些参数的传递 -->
        <el-card style="
                max-width: 600px;
                position:absolute;
                text-align: left;
                pointer-events: none;
                background-color: #f0f0f0;
            "
            :style="{
                visibility: displaybox===1 ? 'visible':'hidden',
                top: boxposition.y +'px',
                left: boxposition.x + 'px'
            }"
        >
            <p class="text_item"><b>file name:</b> {{ boxvalue.file_name }}</p>
            <p class="text_item"><b>function name:</b> {{ boxvalue.function_name }}</p>
            <p class="text_item"><b>frequency:</b> {{ boxvalue.frequency }}</p>
        </el-card>
    </div>
</template>
  
<script>
  import * as d3 from 'd3';
  import {ref,inject} from 'vue';
  const displaybox = ref(0);
  const boxposition = ref({
    "x": 0,
    "y": 0
  });
  const boxvalue = ref({
    "file_name": "",
    "function_name": "",
    "frequency": ""
  });
  export default {
    data() {
      return {
        displaybox,boxvalue,boxposition
      };
    },
    mounted() {
      this.buildCFG();
    },
    setup(){
        const detail_val = inject("detail_val");
        return{
            detail_val
        }
    },
    methods: {
        buildCFG() {
        var that=this;
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
            .style("fill", "#999")
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
            .force("charge", d3.forceManyBody().strength(-10))
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
                    that.detail_val["run_time"]=show_diff_time(last_update - start_time);
                    that.detail_val["cycles_done"]=map["cycles_done"];
                    that.detail_val["total_path"]=map["paths_total"];
                    that.detail_val["unique_crashes"]=map["unique_crashes"];
                    that.detail_val["map_density"]=map["bitmap_cvg"];
                    that.detail_val["execs_done"]=map["execs_done"];
                    that.detail_val["max_depth"]=map["max_depth"];
                    that.detail_val["unique_hangs"]=map["unique_hangs"];

                }).catch(function(error3) {
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
                            search(queue,functions[new_node]["calls"][child]);
                        }
                        for (child in functions[new_node]["refs"]) {
                            // 对search函数进行了修改
                            search(queue,functions[new_node]["refs"][child]);
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
                    .attr('stroke','#999')
                    .attr('stroke-opacity',0.6)
                // .attr("stroke-width", function (d) { return Math.sqrt(d.value); });

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
                            return "purple";
                        }
                        else {
                            return "white";
                        }
                    })
                    .attr("stroke-width", "3px")
                    .attr("fill", function (d) { return d3.interpolatePiYG(1 - d["color"]); })
                    .on("mouseover", function (event, d) {  // 注意这里的参数变化，现在事件对象作为第一个参数
                        // 主要进行参数传递
                        console.log("in mouseover",event,d);
                        that.displaybox=1;
                        // box显示位置
                        that.boxposition.x=event.x + 10;
                        that.boxposition.y=event.y + 28;
                        // box显示内容
                        that.boxvalue.file_name=d["filename"];
                        that.boxvalue.function_name=d["name"];
                        that.boxvalue.frequency=d["freq"];
                        console.log("boxvalue",that.boxvalue)
                    })
                    .on("mouseout", function () {
                        that.displaybox=0;
                    })
                    .on("click", function () {
                        console.log("in click");
                    })
                    .call(d3.drag()
                        .on("start", dragstarted)
                        .on("drag", dragged)
                        .on("end", dragended));
                console.log("circles",circles)
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
                function search(queue,child) {
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

            }).catch(function(error2) {
                console.error("Error loading /api/out/block_freq.json?_=", error2);
            });
        }).catch(function(error) {
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
.text_item{
    margin: 0;
    padding: 0;
    font-size: 16px;
}
</style>
