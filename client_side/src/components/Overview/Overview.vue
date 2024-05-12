<template>
    <div style="height:100%;display:flex;flex-direction:row">
      <div class="chart-container" style="flex:2;height:100%;">
        <el-card class="info-box">
            <h2 class="info-title" style="margin-top:0">Chart</h2>
            <div class="chart-wrapper">
                <canvas id="pathChart"></canvas>
            </div>
        </el-card>
      </div>
      <div class="info-container" style="flex:1;height:100%;">
        <el-card class="info-box">
          <h2 class="info-title" style="margin-top:0">Information</h2>
          <ul class="info-list" id="pathInfo" style="padding:20px;">
            <li><span class="info-label">Run Time:</span> <span class="info-value">{{ chartData.run_time }}</span></li>
            <li><span class="info-label">Cycles Done:</span> <span class="info-value">{{ chartData.cycles_done }}</span></li>
            <li><span class="info-label">Last New Path:</span> <span class="info-value">{{ chartData.last_new_path }}</span></li>
            <li><span class="info-label">Total Path:</span> <span class="info-value">{{ chartData.paths_total }}</span></li>
            <li><span class="info-label">Last Unique Crash:</span> <span class="info-value">{{ chartData.last_unique_crash }}</span></li>
            <li><span class="info-label">Unique Crashes:</span> <span class="info-value">{{ chartData.unique_crashes }}</span></li>
            <li><span class="info-label">Execs per Second:</span> <span class="info-value">{{ chartData.execs_per_sec }}</span></li>
            <li><span class="info-label">Map Density:</span> <span class="info-value">{{ chartData.map_density }}</span></li>
            <li><span class="info-label">Execs Done:</span> <span class="info-value">{{ chartData.execs_done }}</span></li>
            <li><span class="info-label">Stability:</span> <span class="info-value">{{ chartData.stability }}</span></li>
            <li><span class="info-label">Execs Timeout:</span> <span class="info-value">{{ chartData.exec_timeout }}</span></li>
            <li><span class="info-label">Max Depth:</span> <span class="info-value">{{ chartData.max_depth }}</span></li>
            <li><span class="info-label">Command Line:</span> <span class="info-value">{{ chartData.command_line }}</span></li>
          </ul>
        </el-card>
    </div>
    </div>
  </template>
  
  <script>
  import { ref, onMounted, onBeforeUnmount } from 'vue';
  import Chart from 'chart.js/auto';
  import axios from 'axios';
  
  export default {
    name: 'OVERVIEW',
    setup() {
      const chartData = ref({
        run_time: '00:00:00',
        cycles_done: 0,
        last_path: 0,
        paths_total: 0,
        last_crash: 0,
        unique_crashes: 0,
        execs_per_sec: 0,
        bitmap_cvg: 0,
        execs_done: 0,
        stability: '100.00%',
        exec_timeout: 20,
        max_depth: 4,
        command_line: '',
        // command_line: '$TOOL_PATH/fuzz/afl/afl-fuzz -i in -o out ./app @@',
        labels: [],
        datasets: [
          {
            label: 'Paths Total',
            data: [], // Y data
            borderColor: '#3e95cd',
            backgroundColor: 'rgba(75, 192, 192, 0.2)',
            borderWidth: 2,
            pointRadius: 0, // hide points
            fill: false,
            // lineTension: 0.5, // more smooth
          },
        ],
      });
  
      let lineChart = null;
      let timer = null;
      // make chart
      const initializeChart = () => {
        const canvas = document.getElementById('pathChart');
        // canvas的高度填充祖父元素
        canvas.height = canvas.parentElement.parentElement.clientHeight-10;
        if (!canvas) {
          console.error('Canvas element not found.');
          return;
        }
        const ctx = canvas.getContext('2d');
        lineChart = new Chart(ctx, {
          type: 'line',
          data: chartData.value,
          options: {
            animation: false,
            elements: {
              line: {
                tension: 0, // Disable bezier curve
              },
            },
            legend: {
              display: false,
            },
            scales: {
              y: {
                scaleLabel: {
                  display: true,
                  labelString: 'Paths Total',
                },
                grid: {
                  display: true, 
                },
              },
              x: {
                scaleLabel: {
                  display: true,
                  labelString: 'Time',
                  
                },
                ticks: {
                  display: true,
                  callback: function(value) {
                    return get_time(value);
                  }
                },
                grid: {
                    borderDash: [5, 10]  // 设置虚线样式
                }
              },
            },
          },
        });
      };
  
      // chart data : fetch and update
      const fetchDataAndUpdateChart = async () => {
            try {
                const response = await axios.get('/api/out/plot_data?_=' + Math.random());
                const data = response.data.split('\n').slice(1); // remove old line
                const pathsTotal = data.map(line => {
                    const values = line.trim().split(/\s+/); 
                    return parseFloat(values[3]);
                });
        
                // update chart data
                chartData.value.datasets[0].data = pathsTotal;
                chartData.value.labels = Array.from(Array(pathsTotal.length).keys()).map(()=>'');
                lineChart.update();
            } catch (error) {
                console.error('Error fetching and updating data:', error);
            }
      };
      
      // time : xx hour xx min xx sec
      function show_diff_time(startTime, lastUpdate) {
        const diffSeconds = lastUpdate - startTime;
        const hours = Math.floor(diffSeconds / 3600);
        const minutes = Math.floor((diffSeconds % 3600) / 60);
        const seconds = diffSeconds % 60;
        return `${hours} hour ${minutes} min ${seconds} sec`;
      }
      // get : xx hour xx min xx sec
      function get_time(time){
        const hours = Math.floor(time / 3600);
        const minutes = Math.floor((time % 3600) / 60);
        // const seconds = time % 60;
        return `${hours} hour ${minutes} min`;
      }
      // information data : fetch and update
      const fetchDataAndUpdateInformation = async () => {
        try {
          const response = await axios.get("/api/out/fuzzer_stats?_=" + Math.random());
          const text = response.data;
          const lines = text.split('\n');
          const updatedData = {};
          for (const line of lines) {
            const [key, value] = line.split(':').map((str) => str.trim());
            updatedData[key] = value;
          }

          // Calculate last unique crash
          if (updatedData.last_crash !== undefined && parseInt(updatedData.last_crash) !== 0 && updatedData.last_update !== undefined) {
              let lastCrashTime = parseInt(updatedData.last_crash);
              let lastUpdateTime = parseInt(updatedData.last_update);
              let lastUniqueCrash = show_diff_time(lastCrashTime, lastUpdateTime);
              chartData.value.last_unique_crash = lastUniqueCrash;
          } else {
              chartData.value.last_unique_crash = "not yet";
          }
          
          chartData.value.run_time = show_diff_time(parseInt(updatedData.start_time), parseInt(updatedData.last_update)) || chartData.value.run_time;
          chartData.value.last_new_path = show_diff_time(parseInt(updatedData.last_path), parseInt(updatedData.last_update),) || chartData.value.run_time;

          // chartData.value.last_path = show_diff_time(parseInt(updatedData.last_path) || chartData.value.last_path);
          // chartData.value.last_crash = show_diff_time(parseInt(updatedData.last_crash) || chartData.value.last_crash);
          // chartData.value.run_time = updatedData.run_time || chartData.value.run_time;
          // chartData.value.run_time = updatedData.start_time || chartData.value.start_time;
          chartData.value.cycles_done = parseInt(updatedData.cycles_done) || chartData.value.cycles_done;
          // chartData.value.last_path = parseInt(updatedData.last_path) || chartData.value.last_path;
          // chartData.value.last_crash = parseInt(updatedData.last_crash) || chartData.value.last_crash;
          chartData.value.unique_crashes = parseInt(updatedData.unique_crashes) || chartData.value.unique_crashes;
          chartData.value.paths_total = parseInt(updatedData.paths_total) || chartData.value.paths_total;
          chartData.value.execs_per_sec = parseFloat(updatedData.execs_per_sec) || chartData.value.execs_per_sec;
          chartData.value.map_density = parseFloat(updatedData.bitmap_cvg) || chartData.value.bitmap_cvg;
          chartData.value.exec_done = parseInt(updatedData.execs_done) || chartData.value.exec_done;
          chartData.value.stability = updatedData.stability || chartData.value.stability;
          chartData.value.exec_timeout = parseInt(updatedData.exec_timeout) || chartData.value.exec_timeout;
          chartData.value.max_depth = parseInt(updatedData.max_depth) || chartData.value.max_depth;
          chartData.value.command_line = updatedData.command_line || chartData.value.command_line;
  
            // 设置lineChart的横坐标时间

          lineChart.update();
        } catch (error) {
          console.error('Error fetching and updating data:', error);
        }
      };
  
  
      onMounted(() => {
        initializeChart();
        fetchDataAndUpdateInformation(); 
        fetchDataAndUpdateChart(); 
        timer = setInterval(() => {
          fetchDataAndUpdateInformation(); 
          fetchDataAndUpdateChart();
        }, 5000); // every 5s update
      });
  
      onBeforeUnmount(() => {
        clearInterval(timer);
      });
  
  
      return {
        chartData,
        // command_line: '$TOOL_PATH/fuzz/afl/afl-fuzz -i in -o out ./app @@',
      };
  
    },
  };
  </script>
  
  <style scoped>
  .page-container {
    display: flex;
    height: 100vh;
    justify-content: center; /* horizontally centered */
    align-items: center; /* vertically centered */
  }
  
  .chart-container {
    display: flex;
    flex-direction: column;
    /* justify-content: center; vertically centered */
    align-items: center; /* horizontally centered */
    width: 50%;
    height: 100%;
  }
  
  .chart-wrapper {
    width: 100%;
    height: 100%;
    border: 1px solid #ddd;
    border-radius: 5px;
    overflow: hidden;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }
  
  .chart-title {
    margin-bottom: 10px;
    font-size: 24px;
  }
  
  .info-container {
    display: flex;
    justify-content: flex-start; /* left */
    align-items: flex-start; /* top */
    width: 50%;
    height: 100%;
    padding-left: 20px; /* to left distance */
    /* margin-top: auto;  */
    
  }
  
  .info-box {
    height: 100%;
    width: 100%;
  }
  
  .info-title {
    font-weight: bold;
    text-align: center;
    /* color: #1cd6bd; */
  }
  
  .info-list {
    list-style-type: none;
    border-radius: 5px;
    background-color: #333;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    padding: 0;
    
  }
  
  .info-list li {
    margin-bottom: 15px;
    font-size: 18px;
    line-height: 1.5;
    color: #fff;
    text-align: left;
  }
  
  .info-label {
    color: #1cd6bd;
    font-size: 18px;
    margin-right: 5px;
  }
  
  .info-value {
    color: #fff;
    font-size: 18px;
    margin-left: 0;
    padding: 5px;
    /* background-color: #333;
    border-radius: 5px; */
  }
  </style>
  