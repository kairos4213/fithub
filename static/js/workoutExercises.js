function sortExercises() {
  htmx.onLoad(function (content) {
    var sortables = content.querySelectorAll("#workout-exercises");
    for (var i = 0; i < sortables.length; i++) {
      var sortable = sortables[i];
      var sortableInstance = new Sortable(sortable, {
        animation: 150,
        ghostClass: "blue-background-class",

        onEnd: function () {
          this.option("disabled", true);
        },
      });

      sortable.addEventListener("htmx:afterSwap", function () {
        sortableInstance.option("disabled", false);
      });
    }
  });
}

function workoutExerciseRow(
  plannedSets = 1,
  plannedReps = [1],
  plannedWeights = [0],
  completedSets = 0,
  completedReps = [0],
  completedWeights = [0],
) {
  return {
    editingExercise: false,
    plannedSets: plannedSets,
    plannedReps: plannedReps,
    plannedWeights: plannedWeights,
    completedSets: completedSets,
    completedReps: completedReps,
    completedWeights: completedWeights,
    updateArrays(array, n) {
      n = Number(n) || 0;
      if (n < 0) {
        n = 0;
      }

      if (array == "plannedSets") {
        if (this.plannedReps.length > n) {
          this.plannedReps.splice(n);
        } else {
          for (let i = this.plannedReps.length; i < n; i++) {
            this.plannedReps.push(1);
          }
        }

        if (this.plannedWeights.length > n) {
          this.plannedWeights.splice(n);
        } else {
          for (let i = this.plannedWeights.length; i < n; i++) {
            this.plannedWeights.push(0);
          }
        }
      } else {
        if (this.completedReps.length > n) {
          this.completedReps.splice(n);
        } else {
          for (let i = this.completedReps.length; i < n; i++) {
            this.completedReps.push(0);
          }
        }

        if (this.completedWeights.length > n) {
          this.completedWeights.splice(n);
        } else {
          for (let i = this.completedWeights.length; i < n; i++) {
            this.completedWeights.push(0);
          }
        }
      }
    },
  };
}

function workoutExerciseFooter() {
  return {
    open: false,
    plannedSets: 1,
    plannedReps: [1],
    plannedWeights: [0],
    updateArrays(n) {
      n = Number(n) || 0;
      if (n < 0) {
        n = 0;
      }

      if (this.plannedReps.length > n) {
        this.plannedReps.splice(n);
      } else {
        for (let i = this.plannedReps.length; i < n; i++) {
          this.plannedReps.push(1);
        }
      }

      if (this.plannedWeights.length > n) {
        this.plannedWeights.splice(n);
      } else {
        for (let i = this.plannedWeights.length; i < n; i++) {
          this.plannedWeights.push(0);
        }
      }
    },
  };
}
